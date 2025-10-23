package approvaltype

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	investmentEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	benefitparticipation "github.com/raymondsugiarto/coffee-api/pkg/module/benefit_participation"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	"gorm.io/gorm"
)

type benefitParticipationService struct {
	benefitParticipationSvc benefitparticipation.Service
	investmentSvc           investment.Service
	investmentPaymentSvc    investmentpayment.Service
}

func NewBenefitParticipationService(benefitParticipationSvc benefitparticipation.Service, investmentSvc investment.Service, investmentPaymentSvc investmentpayment.Service) *benefitParticipationService {
	return &benefitParticipationService{
		benefitParticipationSvc: benefitParticipationSvc,
		investmentSvc:           investmentSvc,
		investmentPaymentSvc:    investmentPaymentSvc,
	}
}

func (s *benefitParticipationService) ConfirmationApprovalCallback(ctx context.Context, req *entity.ApprovalDto, tx *gorm.DB) (context.Context, error) {
	statusMap := map[model.ApprovalStatus]model.BenefitParticipationStatus{
		model.ApprovalStatusApproved: model.BenefitParticipationStatusActive,
		model.ApprovalStatusRejected: model.BenefitParticipationStatusRejected,
	}

	paymentStatusMap := map[model.ApprovalStatus]model.InvestmentPaymentStatus{
		model.ApprovalStatusApproved: model.InvestmentPaymentStatusSuccess,
		model.ApprovalStatusRejected: model.InvestmentPaymentStatusRejected,
	}

	benefitParticipation, err := s.benefitParticipationSvc.FindByInvestmentPaymentID(ctx, req.RefID)
	if err != nil {
		return ctx, err
	}

	s.benefitParticipationSvc.ConfirmationApprovalCallback(ctx, benefitParticipation.ID, statusMap[model.ApprovalStatus(req.Status)], tx)

	return s.investmentSvc.ConfirmationApprovalCallback(ctx, &investmentEntity.InvestmentPaymentDto{
		ID:     req.RefID,
		Status: paymentStatusMap[model.ApprovalStatus(req.Status)],
	}, tx)
}

func (s *benefitParticipationService) FindByID(ctx context.Context, id string) (interface{}, error) {
	return s.benefitParticipationSvc.FindByInvestmentPaymentID(ctx, id)
}

func (s *benefitParticipationService) NotifyApprovalCallback(ctx context.Context, req *entity.ApprovalDto) error {
	// No implementation needed for benefit participation approval callback notification
	return nil
}
