package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/claim"
	"gorm.io/gorm"
)

type ClaimService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}

type claimService struct {
	claimService claim.Service
}

func NewClaimService(cs claim.Service) ClaimService {
	return &claimService{
		claimService: cs,
	}
}

func (s *claimService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	return s.claimService.ConfirmationApprovalCallback(ctx, &entity.ClaimDto{
		ID:             req.RefID,
		ApprovalStatus: model.ApprovalStatus(req.Status),
	}, tx)
}

func (s *claimService) FindByID(ctx context.Context, id string) (interface{}, error) {
	result, err := s.claimService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *claimService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	// No implementation needed for claim approval callback notification
	return nil
}
