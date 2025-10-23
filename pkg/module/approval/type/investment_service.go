package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	iventity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	"gorm.io/gorm"
)

type InvestmentService interface {
	ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error)
	FindByID(ctx context.Context, id string) (interface{}, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error
}

type investmentService struct {
	investmentService        investment.Service
	investmentPaymentService investmentpayment.Service
}

func NewInvestmentService(is investment.Service, ips investmentpayment.Service) InvestmentService {
	return &investmentService{
		investmentService:        is,
		investmentPaymentService: ips,
	}
}

func (s *investmentService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	paymentStatusMap := map[model.ApprovalStatus]model.InvestmentPaymentStatus{
		model.ApprovalStatusApproved: model.InvestmentPaymentStatusSuccess,
		model.ApprovalStatusRejected: model.InvestmentPaymentStatusRejected,
	}

	return s.investmentService.ConfirmationApprovalCallback(ctx, &iventity.InvestmentPaymentDto{
		ID:     req.RefID,
		Status: paymentStatusMap[model.ApprovalStatus(req.Status)],
	}, tx)
}

func (s *investmentService) FindByID(ctx context.Context, id string) (interface{}, error) {
	return s.investmentPaymentService.FindByID(ctx, id)
}

func (s *investmentService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	paymentStatusMap := map[model.ApprovalStatus]model.InvestmentPaymentStatus{
		model.ApprovalStatusApproved: model.InvestmentPaymentStatusSuccess,
		model.ApprovalStatusRejected: model.InvestmentPaymentStatusRejected,
	}

	return s.investmentService.NotifyApprovalCallback(ctx, &iventity.InvestmentPaymentDto{
		ID:     req.RefID,
		Status: paymentStatusMap[model.ApprovalStatus(req.Status)],
	})
}
