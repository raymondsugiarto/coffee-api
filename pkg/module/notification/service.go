package notification

import (
	"context"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/coffee-api/config"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	ei "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	SendEmailTemplate(ctx context.Context, req *entity.NotificationInputDto) (b.CreateSmtpEmail, error)

	Create(ctx context.Context, dto *entity.NotificationDto) (*entity.NotificationDto, error)
	FindByID(ctx context.Context, id string) (*entity.NotificationDto, error)
	Update(ctx context.Context, dto *entity.NotificationDto) (*entity.NotificationDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *e.NotificationFindAllRequest) (*pagination.ResultPagination, error)

	NotifyInvestment(ctx context.Context, dto *ei.InvestmentDto) error
	NotifyInvestmentPaymentConfirmed(ctx context.Context, dto *ei.InvestmentDto) error
}

type service struct {
	client *b.APIClient
	repo   Repository
}

func NewService(apiClient *b.APIClient, repo Repository) Service {
	return &service{
		client: apiClient,
		repo:   repo,
	}
}

func (s *service) SendEmailTemplate(ctx context.Context, req *entity.NotificationInputDto) (b.CreateSmtpEmail, error) {
	log.Infof("Sending email to %+v", req.To)
	sender := req.From
	if req.From == nil {
		cfg := config.GetConfig()
		sender = &b.SendSmtpEmailSender{
			Email: cfg.Mail.Brevo.Sender,
		}
	}
	request := b.SendSmtpEmail{
		Sender:     sender,
		To:         req.To,
		TemplateId: req.TemplateID,
		Params:     req.Data,
	}
	log.Infof("Sending email to %+v with %+v", req.To, request)
	resp, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(ctx, request)
	log.Infof("Email sent to %+v with response %+v", req.To, resp)
	if err != nil {
		log.Errorf("Failed to send email: %v", err)
		return resp, err
	}
	return resp, nil
}

func (s *service) Create(ctx context.Context, dto *entity.NotificationDto) (*entity.NotificationDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.NotificationDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.NotificationDto) (*entity.NotificationDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.NotificationFindAllRequest) (*pagination.ResultPagination, error) {
	// update read at if not set
	s.repo.UpdateRead(ctx, req.UserID)

	return s.repo.FindAll(ctx, req)
}

func (s *service) NotifyInvestment(ctx context.Context, dto *ei.InvestmentDto) error {
	for _, item := range dto.InvestmentItems {
		notification := &entity.NotificationDto{
			OrganizationID: dto.OrganizationID,
			UserID:         item.CustomerID,
			RefModule:      "investment",
			RefTable:       "investments",
			RefID:          dto.ID,
			RefCode:        dto.Code,
			Description:    "Investment created successfully",
			NotifyAt:       time.Now(),
		}
		if _, err := s.Create(ctx, notification); err != nil {
			log.Errorf("Failed to create notification for investment %s: %v", dto.ID, err)
			return err
		}
	}
	return nil
}

func (s *service) NotifyInvestmentPaymentConfirmed(ctx context.Context, dto *ei.InvestmentDto) error {
	for _, item := range dto.InvestmentItems {
		notification := &entity.NotificationDto{
			OrganizationID: dto.OrganizationID,
			UserID:         item.CustomerID,
			RefModule:      "investment_payment",
			RefTable:       "investment",
			RefID:          dto.ID,
			RefCode:        dto.Code,
			Description:    "Your investment payment has been confirmed successfully",
			NotifyAt:       time.Now(),
		}
		if _, err := s.Create(ctx, notification); err != nil {
			log.Errorf("Failed to create notification for investment %s: %v", dto.ID, err)
			return err
		}
	}
	return nil
}
