package useridentityverification

import (
	"context"
	"errors"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2/log"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	"github.com/raymondsugiarto/coffee-api/pkg/module/whatsapp"
	"github.com/raymondsugiarto/coffee-api/pkg/shared"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

type Callback func(ctx context.Context, dto *entity.UserIdentityVerificationDto) error

type Service interface {
	Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error)
	FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error)
	Verify(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error)
	Resend(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error)
}

type service struct {
	repository          Repository
	callbackMap         map[string]Callback
	customerService     customer.Service
	whatsappService     whatsapp.Service
	notificationService notification.Service
	companyService      company.Service
}

func NewService(
	repository Repository,
	customerService customer.Service,
	whatsappService whatsapp.Service,
	notificationService notification.Service,
	companyService company.Service,
) Service {
	callbackMap := map[string]Callback{
		"FORGOT_PASSWORD_CUSTOMER":      customerService.ChangePassword,
		"SIGN_IN_VERIFY_EMAIL_CUSTOMER": customerService.VerifyEmail,
		"FORGOT_PASSWORD_COMPANY":       companyService.ChangePassword,
	}
	s := &service{
		repository:          repository,
		customerService:     customerService,
		callbackMap:         callbackMap,
		whatsappService:     whatsappService,
		notificationService: notificationService,
		companyService:      companyService,
	}
	return s
}

func (s *service) Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error) {
	dto.TryCount = 0

	response, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error) {
	return s.repository.FindByIDAndUniqueCode(ctx, id, uniqueCode)
}

func (s *service) Verify(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error) {
	userIdentityVerificationDto, err := s.repository.FindByIDAndUniqueCode(ctx, dto.ID, dto.UniqueCode)
	if err != nil {
		return nil, err
	}
	if userIdentityVerificationDto.TryCount >= 3 {
		return nil, status.New(status.BadRequest, errors.New("try count exceeded"))
	}
	if userIdentityVerificationDto.Status == "SUCCESS" {
		return nil, status.New(status.BadRequest, errors.New("already verified"))
	}
	if userIdentityVerificationDto.ExpiredAt.Before(time.Now()) {
		return nil, status.New(status.BadRequest, errors.New("expired"))
	}

	userIdentityVerificationDto.Status = "SUCCESS"

	dto.ID = userIdentityVerificationDto.ID
	dto.UserID = userIdentityVerificationDto.UserID
	dto.UserIdentity = userIdentityVerificationDto.UserIdentity
	dto.IdentityType = userIdentityVerificationDto.IdentityType
	err = s.callbackMap[userIdentityVerificationDto.IdentityFor](ctx, dto)
	if err != nil {
		log.Errorf("errorCallback: %v", err)
		return nil, err
	}

	go s.repository.Update(ctx, userIdentityVerificationDto)

	return dto, nil
}

func (s *service) Resend(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error) {
	userIdentityVerificationDto, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if userIdentityVerificationDto.TryCount >= 3 {
		return nil, status.New(status.BadRequest, errors.New("try count exceeded"))
	}
	if userIdentityVerificationDto.Status == "SUCCESS" {
		return nil, status.New(status.BadRequest, errors.New("already verified"))
	}
	if userIdentityVerificationDto.ExpiredAt.Before(time.Now()) {
		return nil, status.New(status.BadRequest, errors.New("expired"))
	}

	userIdentityVerificationDto.TryCount++
	go s.repository.Update(ctx, userIdentityVerificationDto)

	if userIdentityVerificationDto.IdentityType == "PHONE_NUMBER" {
		// send OTP again
		message := "Your verification code is " + userIdentityVerificationDto.UniqueCode + ". For security reasons, please do not share this code with anyone."
		s.whatsappService.SendMessage(ctx, userIdentityVerificationDto.UserIdentity, message)
	} else if userIdentityVerificationDto.IdentityType == "EMAIL" {

		if userIdentityVerificationDto.IdentityFor == "SIGN_IN_VERIFY_EMAIL_CUSTOMER" {
			err = s.sendEmailVerificationOtpToCustomer(ctx, userIdentityVerificationDto)
			if err != nil {
				log.Errorf("errorSendEmailVerificationOtpToCustomer: %v", err)
				return nil, err
			}
		} else if userIdentityVerificationDto.IdentityFor == "FORGOT_PASSWORD_CUSTOMER" {
			err = s.sendForgotPasswordOtpToCustomer(ctx, userIdentityVerificationDto)
			if err != nil {
				log.Errorf("errorSendForgotPasswordOtpToCustomer: %v", err)
				return nil, err
			}
		} else if userIdentityVerificationDto.IdentityFor == "FORGOT_PASSWORD_COMPANY" {
			err = s.sendForgotPasswordOtpToCompany(ctx, userIdentityVerificationDto)
			if err != nil {
				log.Errorf("errorSendForgotPasswordOtpToCompany: %v", err)
				return nil, err
			}
		} else {
			return nil, status.New(status.BadRequest, errors.New("invalid identity for"))
		}
	}

	return userIdentityVerificationDto, nil
}

func (s *service) sendEmailVerificationOtpToCustomer(ctx context.Context, userIdentityVerificationDto *entity.UserIdentityVerificationDto) error {
	// send email again
	customerDto, err := s.customerService.FindByUserID(ctx, userIdentityVerificationDto.UserID)
	if err != nil {
		log.Errorf("errorFindByUserID: %v", err)
		return err
	}
	s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: shared.TEMPLATE_ID_EMAIL_VERIFICATION,
		To: []b.SendSmtpEmailTo{
			{
				Email: customerDto.Email,
			},
		},
		Data: map[string]interface{}{
			"name":  customerDto.FirstName,
			"token": userIdentityVerificationDto.UniqueCode,
		},
	})
	return nil
}

func (s *service) sendForgotPasswordOtpToCustomer(ctx context.Context, userIdentityVerificationDto *entity.UserIdentityVerificationDto) error {
	// send email again
	customerDto, err := s.customerService.FindByUserID(ctx, userIdentityVerificationDto.UserID)
	if err != nil {
		log.Errorf("errorFindByUserID: %v", err)
		return err
	}
	s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: shared.TEMPLATE_ID_FORGOT_PASSWORD,
		To: []b.SendSmtpEmailTo{
			{
				Email: userIdentityVerificationDto.UserIdentity,
			},
		},
		Data: map[string]interface{}{
			"name":  customerDto.FirstName,
			"token": userIdentityVerificationDto.UniqueCode,
		},
	})
	return nil
}

func (s *service) sendForgotPasswordOtpToCompany(ctx context.Context, userIdentityVerificationDto *entity.UserIdentityVerificationDto) error {
	// send email again
	companyDto, err := s.companyService.FindByUserID(ctx, userIdentityVerificationDto.UserID)
	if err != nil {
		log.Errorf("errorFindByUserID: %v", err)
		return err
	}
	s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: shared.TEMPLATE_ID_FORGOT_PASSWORD,
		To: []b.SendSmtpEmailTo{
			{
				Email: companyDto.Email,
			},
		},
		Data: map[string]interface{}{
			"name":  companyDto.FirstName,
			"token": userIdentityVerificationDto.UniqueCode,
		},
	})
	return nil
}
