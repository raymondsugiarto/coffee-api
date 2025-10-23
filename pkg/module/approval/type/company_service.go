package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"gorm.io/gorm"
)

type CompanyService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}

type companyService struct {
	companyService company.Service
}

func NewCompanyService(cs company.Service) CompanyService {
	return &companyService{
		companyService: cs,
	}
}

func (s *companyService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	return s.companyService.ConfirmationApprovalCallback(ctx, &entity.CompanyDto{
		// TODO: map here
		ID:     req.RefID,
		Status: req.Status,
	}, tx)
}

func (s *companyService) FindByID(ctx context.Context, id string) (interface{}, error) {
	result, err := s.companyService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *companyService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	// No implementation needed for company approval callback notification
	return nil
}
