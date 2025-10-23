package approval

import (
	"context"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	sh "github.com/raymondsugiarto/coffee-api/pkg/shared"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ApprovalDto) (*entity.ApprovalDto, error)
	FindByID(ctx context.Context, id string) (*entity.ApprovalDto, error)
	FindByRefID(ctx context.Context, refId, approvalType string) (*entity.ApprovalDto, error)
	Update(ctx context.Context, dto *entity.ApprovalDto) (*entity.ApprovalDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ApprovalFindAllRequest) (*pagination.ResultPagination, error)
	Confirmation(ctx context.Context, req *entity.ApprovalDto) (*entity.ApprovalDto, error)
	CreateWithTx(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error)
	UpdateWithTx(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error)

	SetCallbackCompany(companyService CallbackService)
	SetCallbackCustomer(customerService CallbackService)
	SetCallbackClaim(claimService CallbackService)
	SetCallbackInvestment(investmentService CallbackService)
	SetCallbackTicket(ticketService CallbackService)
	SetCallbackBenefitParticipation(benefitParticipationService CallbackService)
}

type service struct {
	repo                Repository
	mapCallback         map[model.ApprovalType]CallbackService
	notificationService notification.Service
}

func NewService(
	repo Repository,
	notificationService notification.Service,
) Service {
	mc := make(map[model.ApprovalType]CallbackService)
	return &service{
		repo,
		mc,
		notificationService,
	}
}

func (s *service) SetCallbackCompany(companyService CallbackService) {
	s.mapCallback[model.ApprovalTypeCompany] = companyService
}

func (s *service) SetCallbackCustomer(customerService CallbackService) {
	s.mapCallback[model.ApprovalTypeCustomer] = customerService
}

func (s *service) SetCallbackClaim(claimService CallbackService) {
	s.mapCallback[model.ApprovalTypeClaim] = claimService
}

func (s *service) SetCallbackInvestment(investmentService CallbackService) {
	s.mapCallback[model.ApprovalTypeInvestment] = investmentService
}

func (s *service) SetCallbackTicket(ticketService CallbackService) {
	s.mapCallback[model.ApprovalTicketInvestment] = ticketService
}

func (s *service) SetCallbackBenefitParticipation(benefitParticipationService CallbackService) {
	s.mapCallback[model.ApprovalTypeBenefitParticipation] = benefitParticipationService
}

func (s *service) Confirmation(ctx context.Context, req *entity.ApprovalDto) (*entity.ApprovalDto, error) {
	approval, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	approval.Status = req.Status

	res, err := s.repo.Confirmation(ctx, approval, func(tx *gorm.DB) error {
		if ctx, err = s.mapCallback[approval.Type].ConfirmationApprovalCallback(ctx, approval, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	go s.mapCallback[approval.Type].NotifyApprovalCallback(ctx, approval)

	if req.Status == "REJECTED" || req.Status == string(model.ApprovalStatusRejected) {
		s.sendRejectedEmail(ctx, approval)
	}
	return res, nil
}

func (s *service) Create(ctx context.Context, dto *entity.ApprovalDto) (*entity.ApprovalDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.CreatedDate = time.Now().Local().UTC()
	return s.repo.Create(ctx, dto, nil)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ApprovalDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) FindByRefID(ctx context.Context, refID, approvalType string) (*entity.ApprovalDto, error) {
	return s.repo.GetByRefID(ctx, refID, approvalType)
}

func (s *service) Update(ctx context.Context, dto *entity.ApprovalDto) (*entity.ApprovalDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto, nil)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.ApprovalFindAllRequest) (*pagination.ResultPagination, error) {
	res, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	results := res.Data.([]*entity.ApprovalDto)
	for i, item := range results {
		if callback, ok := s.mapCallback[item.Type]; ok {
			refData, err := callback.FindByID(ctx, item.RefID)
			if err == nil && refData != nil {
				results[i].RefData = refData
			}
		}
	}
	res.Data = results

	return res, nil
}

func (s *service) CreateWithTx(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.CreatedDate = time.Now().Local().UTC()
	return s.repo.Create(ctx, dto, tx)
}

func (s *service) UpdateWithTx(ctx context.Context, dto *entity.ApprovalDto, tx *gorm.DB) (*entity.ApprovalDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto, tx)
}

func (s *service) sendRejectedEmail(ctx context.Context, approval *entity.ApprovalDto) {
	var email, name, description string

	if approval.RefData == nil {
		if callback, ok := s.mapCallback[approval.Type]; ok {
			if refData, err := callback.FindByID(ctx, approval.RefID); err == nil {
				approval.RefData = refData
			}
		}
	}
	if contact, ok := approval.RefData.(entity.HasContactInfo); ok {
		info := contact.GetInfo()
		email = info.Email
		name = info.Name
		description = info.Description
	}

	go s.notificationService.SendEmailTemplate(ctx, &entity.NotificationInputDto{
		TemplateID: sh.REJECTED_EMAIL_TEMPLATE_ID,
		To: []b.SendSmtpEmailTo{
			{Email: email},
		},
		Data: map[string]interface{}{
			"name":        name,
			"description": description,
		},
	})
}
