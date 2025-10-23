package participant

import (
	"context"
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ParticipantDto, customer *entity.CustomerDto) (*entity.ParticipantDto, error)
	FindByID(ctx context.Context, id string) (*entity.ParticipantDto, error)
	Update(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error)
	UpdateStatus(ctx context.Context, id string, status model.ParticipantStatus) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ParticipantFindAllRequest) (*pagination.ResultPagination, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ParticipantDto, error)

	FindAllParticipantCompany(ctx context.Context, req *entity.ParticipantFindAllRequest) ([]*entity.ParticipantDto, error)
}

type service struct {
	repo            Repository
	unitLinkService unitlink.Service
}

func NewService(repo Repository, unitLinkService unitlink.Service) Service {
	return &service{repo, unitLinkService}
}

// ketika service register employee, pasti type DKP
func (s *service) Create(ctx context.Context, dto *entity.ParticipantDto, customer *entity.CustomerDto) (*entity.ParticipantDto, error) {
	// Determine if this is a group participant by checking customer's CompanyID
	isGroup := false
	if customer != nil && customer.CompanyID != "" {
		isGroup = true
	}

	// Generate new participant number format
	participantNumber, err := s.generateParticipantNumber(isGroup, time.Now())
	if err != nil {
		return nil, err
	}
	dto.Code = participantNumber

	if dto.Status == "" {
		dto.Status = model.ParticipantStatusInactive
	}
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ParticipantDto, error) {
	dto, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	sumUnitLink, err := s.unitLinkService.SumInvestmentProductByParticipant(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	dto.Balance = sumUnitLink.TotalAmount
	return dto, nil
}

func (s *service) Update(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error) {
	return s.repo.Update(ctx, dto)
}

func (s *service) UpdateStatus(ctx context.Context, id string, status model.ParticipantStatus) error {
	participant, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if participant.Status == status {
		return nil
	}

	participant.Status = status
	_, err = s.repo.Update(ctx, participant)
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.ParticipantFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) FindAllParticipantCompany(ctx context.Context, req *entity.ParticipantFindAllRequest) ([]*entity.ParticipantDto, error) {
	companyID := shared.GetCompanyID(ctx)
	req.CompanyID = companyID
	return s.repo.FindAllParticipantCompany(ctx, req)
}

func (s *service) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ParticipantDto, error) {
	return s.repo.FindByCompanyID(ctx, companyID)
}

// generateParticipantNumber generates a new participant number based on type and creation date
// Format: SIM-IND-YYYYMM-XXXXXX (individual) or SIM-INS-YYYYMM-XXXXXX (group)
func (s *service) generateParticipantNumber(isGroup bool, createdAt time.Time) (string, error) {
	// Determine participant type prefix
	participantType := "IND"
	if isGroup {
		participantType = "INS"
	}

	yearMonth := createdAt.Format("200601")

	// Generate random 6-digit number
	randomSuffix, err := gonanoid.Generate("0123456789", 6)
	if err != nil {
		return "", err
	}

	// Combine all parts
	return fmt.Sprintf("SIM-%s-%s-%s", participantType, yearMonth, randomSuffix), nil
}
