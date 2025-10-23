package claim

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ClaimDto) (*entity.ClaimDto, error)
	ConfirmationApprovalCallback(ctx context.Context, dto *entity.ClaimDto, tx *gorm.DB) (context.Context, error)
	FindAll(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error)
	FindAllByCompany(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.ClaimDto, error)

	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ClaimDto, error)
}

type service struct {
	repository         Repository
	approvalService    approval.Service
	unitLinkSvc        unitlink.Service
	participantService participant.Service
}

func NewService(
	repo Repository, approvalSvc approval.Service, unitLinkSvc unitlink.Service,
	participantService participant.Service,
) Service {
	return &service{
		repository:         repo,
		approvalService:    approvalSvc,
		unitLinkSvc:        unitLinkSvc,
		participantService: participantService,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.ClaimDto) (*entity.ClaimDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	uid := shared.GetUserCredential(ctx).UserID
	dto.ApprovalStatus = model.ApprovalStatusSubmit // Default to SUBMIT

	if dto.OrganizationID == "" {
		dto.OrganizationID = organizationID
	}

	participant, err := s.participantService.FindByID(ctx, dto.ParticipantID)
	if err != nil {
		return nil, status.New(status.EntityNotFound, errors.New("participant not found"))
	}
	dto.Participant = participant

	return s.repository.Create(ctx, dto, func(tx *gorm.DB) error {
		if _, err := s.approvalService.CreateWithTx(ctx, dto.ToApprovalSubmitDto(uid), tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *service) ConfirmationApprovalCallback(ctx context.Context, dto *entity.ClaimDto, tx *gorm.DB) (context.Context, error) {
	claim, err := s.repository.FindByID(ctx, dto.ID)
	if err != nil {
		return ctx, err
	}

	claim.ApprovalStatus = dto.ApprovalStatus

	if _, err := s.repository.Update(ctx, claim, tx); err != nil {
		return ctx, err
	}

	if claim.ApprovalStatus != model.ApprovalStatusApproved {
		return ctx, nil
	}

	err = s.unitLinkSvc.ClaimUnitLink(ctx, tx, claim.Participant.ID)
	if err != nil {
		return ctx, err
	}

	return ctx, err
}

func (s *service) FindAll(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindAllByCompany(ctx context.Context, req *entity.ClaimFindAllRequest) (*pagination.ResultPagination, error) {
	companyID := shared.GetCompanyID(ctx)
	if companyID == nil {
		return nil, status.New(status.EntityNotFound, errors.New("companyIDNotFound"))
	}
	req.CompanyID = companyID
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ClaimDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ClaimDto, error) {
	return s.repository.FindByCompanyID(ctx, companyID)
}
