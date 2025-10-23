package authentication

import (
	"context"
	"errors"
	"fmt"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication/token"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	customerpoint "github.com/raymondsugiarto/coffee-api/pkg/module/customer/customer_point"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
	useridentityverification "github.com/raymondsugiarto/coffee-api/pkg/module/user_identity_verification"
	"github.com/raymondsugiarto/coffee-api/pkg/module/whatsapp"
	sh "github.com/raymondsugiarto/coffee-api/pkg/shared"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
	"gorm.io/gorm"
)

type Service interface {
	SignUpCustomer(context.Context, *entity.SignUpCustomerDto) (*entity.SignUpCustomerDto, error)
	SignUpCompany(context.Context, *entity.SignUpCompanyDto) (*entity.SignUpCompanyDto, error)
	SignIn(context.Context, *entity.LoginRequestDto) (*entity.LoginDto, error)
	ForgotPasswordCustomer(context.Context, *entity.ForgotPasswordDto) (*entity.ForgotPasswordDto, error)
	ForgotPasswordCompany(ctx context.Context, request *entity.ForgotPasswordDto) (*entity.ForgotPasswordDto, error)
}

type service struct {
	userCredentialService           usercredential.Service
	customerService                 customer.Service
	userIdentityVerificationService useridentityverification.Service
	notificationService             notification.Service
	companyService                  company.Service
	whatsappService                 whatsapp.Service
	tokenService                    token.Service
	customerPointService            customerpoint.Service
	adminService                    admin.Service
}

func NewService(
	userCredentialService usercredential.Service,
	customerService customer.Service,
	userIdentityVerificationService useridentityverification.Service,
	notificationService notification.Service,
	companyService company.Service,
	whatsappService whatsapp.Service,
	tokenService token.Service,
	customerPointService customerpoint.Service,
	adminService admin.Service,
) Service {
	return &service{
		userCredentialService:           userCredentialService,
		customerService:                 customerService,
		userIdentityVerificationService: userIdentityVerificationService,
		notificationService:             notificationService,
		companyService:                  companyService,
		whatsappService:                 whatsappService,
		tokenService:                    tokenService,
		customerPointService:            customerPointService,
		adminService:                    adminService,
	}
}

func (s *service) SignUpCustomer(ctx context.Context, request *entity.SignUpCustomerDto) (*entity.SignUpCustomerDto, error) {
	dto := request.ToCustomerDto(ctx)
	if request.ReferralCode != "" {
		parent, err := s.customerService.FindByReferralCode(ctx, request.ReferralCode)
		if err != nil {
			log.Errorf("errorFindByReferralCode: %v", err)
			return nil, status.New(status.BadRequest, err)
		}
		dto.CustomerIDParent = parent.ID
	}

	// TODO: validate username
	uc, err := s.userCredentialService.FindByUsername(ctx, request.Username)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("errorFindByUsername: %v", err)
			return nil, status.New(status.BadRequest, err)
		}
	}

	if uc != nil {
		return nil, status.New(status.BadRequest, errors.New("usernameAlreadyExist"))
	}

	// validate email
	c, err := s.customerService.FindByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("errorFindByEmail: %v", err)
			return nil, status.New(status.BadRequest, err)
		}
	}
	if c != nil {
		return nil, status.New(status.BadRequest, errors.New("emailAlreadyExist"))
	}

	dto.ApprovalStatus = model.ApprovalStatusKyc
	dto.ReferralCode, err = gonanoid.Generate(utils.ALPHA_NUMERIC, 6)
	if err != nil {
		log.Errorf("errorGenerateReferralCode: %v", err)
		return nil, status.New(status.BadRequest, err)
	}
	customerDto, err := s.customerService.Create(ctx, dto)
	if err != nil {
		log.Errorf("errorCreateCustomer: %v", err)
		return nil, status.New(status.BadRequest, err)
	}
	dto.ID = customerDto.ID
	request.ID = customerDto.ID
	request.ReferralCode = dto.ReferralCode

	if dto.CustomerIDParent != "" {
		// add point to parent
		go s.customerPointService.Create(ctx, dto.ToCustomerPointDto())
	}

	return request, nil
}

func (s *service) SignUpCompany(ctx context.Context, request *entity.SignUpCompanyDto) (*entity.SignUpCompanyDto, error) {
	dto := request.ToCompanyDto(ctx)

	// TODO: validate username
	uc, err := s.userCredentialService.FindByUsername(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("errorFindByUsername: %v", err)
			return nil, status.New(status.BadRequest, fmt.Errorf("username %s sudah digunakan", request.Email))
		}
	}

	if uc != nil {
		return nil, status.New(status.BadRequest, fmt.Errorf("username %s sudah terdaftar", request.Email))
	}

	// validate email
	c, err := s.companyService.FindByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("errorFindByEmail: %v", err)
			return nil, status.New(status.BadRequest, fmt.Errorf("email %s sudah digunakan", request.Email))
		}
	}
	if c != nil {
		return nil, status.New(status.BadRequest, fmt.Errorf("email %s sudah terdaftar", request.Email))
	}

	companyDto, err := s.companyService.Create(ctx, dto)
	if err != nil {
		log.Errorf("errorCreateCompany: %v", err)
		return nil, status.New(status.InternalServerError, err)
	}
	request.ID = companyDto.ID
	return request, nil
}

func (s *service) SignIn(ctx context.Context, request *entity.LoginRequestDto) (*entity.LoginDto, error) {
	userCredentialDto, err := s.userCredentialService.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	userCredentialData := e.UserCredentialData{
		ID:     userCredentialDto.ID,
		UserID: userCredentialDto.User.ID,
	}

	if userCredentialDto.User.UserType == "COMPANY" {
		admin, err := s.adminService.FindByUserID(ctx, userCredentialDto.User.ID)
		if err != nil {
			return nil, err
		}
		if admin != nil && admin.CompanyID != "" {
			userCredentialData.CompanyID = admin.CompanyID
		}

	} else if userCredentialDto.User.UserType == "CUSTOMER" {
		cust, err := s.customerService.FindByUserID(ctx, userCredentialDto.User.ID)
		if err != nil {
			return nil, err
		}
		userCredentialData.CustomerID = cust.ID
		if userCredentialDto.User.EmailVerificationStatus == "UNVERIFIED" {
			userIdentityVerificationDto, err := s.sendEmailVerificationOtpToCustomer(ctx, userCredentialDto.User.ID)
			if err != nil {
				log.Errorf("errorSendEmailVerificationOtpToCustomer: %v", err)
				return nil, err
			}
			return &entity.LoginDto{
				Status:                      "UNVERIFIED",
				UserIdentityVerificationDto: userIdentityVerificationDto,
			}, nil
		}
	}

	// pp, _ := utils.HashPassword(request.Password)
	// fmt.Printf("password hash: %+v\n", pp)
	// fmt.Printf("userCredentialData: %+v ::: %+v\n", request.Password, userCredentialDto.Password)
	if !utils.CheckPasswordHash(request.Password, userCredentialDto.Password) {
		return nil, errors.New("invalidPassword")
	}
	return s.tokenService.GenerateToken(ctx, userCredentialData)
}

func (s *service) ForgotPasswordCustomer(ctx context.Context, request *entity.ForgotPasswordDto) (*entity.ForgotPasswordDto, error) {
	customerDto, err := s.customerService.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if customerDto == nil {
		return nil, errors.New("emailNotFound")
	}

	dto := request.ToUserIdentityVerificationDto(model.FORGOT_PASSWORD_CUSTOMER)
	dto.ExpiredAt = time.Now().Add(time.Minute * 60).Local().UTC()
	dto.UserID = customerDto.UserID
	dto.Status = "REQUEST"
	uniqueCode, _ := gonanoid.Generate(utils.NUMERIC, 6)
	dto.UniqueCode = uniqueCode
	userIdentityVerification, err := s.userIdentityVerificationService.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	request.UserIdentityVerificationID = userIdentityVerification.ID

	go s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: sh.TEMPLATE_ID_FORGOT_PASSWORD,
		To: []b.SendSmtpEmailTo{
			{
				Email: customerDto.Email,
			},
		},
		Data: map[string]interface{}{
			"name":  customerDto.FirstName,
			"token": uniqueCode,
		},
	})
	return request, nil
}

func (s *service) sendEmailVerificationOtpToCustomer(ctx context.Context, userID string) (*entity.UserIdentityVerificationDto, error) {
	customer, err := s.customerService.FindByUserID(ctx, userID)
	if err != nil {
		log.Errorf("errorFindByUserID: %v", err)
		return nil, err
	}
	dto := &entity.UserIdentityVerificationDto{
		UserIdentity: customer.Email,
		IdentityFor:  "SIGN_IN_VERIFY_EMAIL_CUSTOMER",
		IdentityType: "EMAIL",
	}
	dto.ExpiredAt = time.Now().Add(time.Minute * 60).Local().UTC()
	dto.UserID = userID
	dto.Status = "REQUEST"
	uniqueCode, _ := gonanoid.Generate(utils.NUMERIC, 5)
	dto.UniqueCode = uniqueCode
	useridentityverification, err := s.userIdentityVerificationService.Create(ctx, dto)
	if err != nil {
		log.Errorf("errorCreateUserIdentityVerification: %v", err)
		return nil, err
	}
	go s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: sh.TEMPLATE_ID_EMAIL_VERIFICATION,
		To: []b.SendSmtpEmailTo{
			{
				Email: customer.Email,
			},
		},
		Data: map[string]interface{}{
			"name":  customer.FirstName,
			"token": uniqueCode,
		},
	})
	return useridentityverification, nil
}

func (s *service) ForgotPasswordCompany(ctx context.Context, request *entity.ForgotPasswordDto) (*entity.ForgotPasswordDto, error) {
	companyDto, err := s.companyService.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if companyDto == nil {
		return nil, errors.New("emailNotFound")
	}

	dto := request.ToUserIdentityVerificationDto(model.FORGOT_PASSWORD_COMPANY)
	dto.ExpiredAt = time.Now().Add(time.Minute * 60).Local().UTC()
	dto.UserID = companyDto.UserID
	dto.Status = "REQUEST"
	uniqueCode, _ := gonanoid.Generate(utils.NUMERIC, 6)
	dto.UniqueCode = uniqueCode
	userIdentityVerification, err := s.userIdentityVerificationService.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	request.UserIdentityVerificationID = userIdentityVerification.ID

	go s.notificationService.SendEmailTemplate(ctx, &e.NotificationInputDto{
		TemplateID: sh.TEMPLATE_ID_FORGOT_PASSWORD,
		To: []b.SendSmtpEmailTo{
			{
				Email: companyDto.Email,
			},
		},
		Data: map[string]interface{}{
			"name":  companyDto.LastName,
			"token": uniqueCode,
		},
	})
	return request, nil
}
