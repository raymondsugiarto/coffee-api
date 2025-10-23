package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"gorm.io/gorm"
)

type CustomerService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}

type customerService struct {
	customerService customer.Service
}

func NewCustomerService(cs customer.Service) CustomerService {
	return &customerService{
		customerService: cs,
	}
}

func (s *customerService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	return s.customerService.ConfirmationApprovalCallback(ctx, &entity.CustomerDto{
		// TODO: map here
		ID:             req.RefID,
		ApprovalStatus: model.ApprovalStatus(req.Status),
	}, tx)
}

func (s *customerService) FindByID(ctx context.Context, id string) (interface{}, error) {
	result, err := s.customerService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *customerService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	// No implementation needed for customer approval callback notification
	return nil
}
