package benefitparticipation

import (
	"context"
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	investmentEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/benefit_type"
	participant "github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	feesetting "github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *entity.BenefitParticipationDto) (*entity.BenefitParticipationDto, error)
	FindByID(ctx context.Context, id string) (*entity.BenefitParticipationDto, error)
	FindByInvestmentPaymentID(ctx context.Context, investmentPaymentID string) (*entity.BenefitParticipationDto, error)
	ConfirmationApprovalCallback(ctx context.Context, benefitParticipationID string, status model.BenefitParticipationStatus, tx *gorm.DB) error
}

type service struct {
	repository                    Repository
	investmentService             investment.Service
	investmentDistributionService investmentdistribution.Service
	participantService            participant.Service
	benefitTypeService            benefit_type.Service
	feeSettingService             feesetting.Service
}

func NewService(
	repository Repository,
	investmentService investment.Service,
	investmentDistributionService investmentdistribution.Service,
	participantService participant.Service,
	benefitTypeService benefit_type.Service,
	feeSettingService feesetting.Service,
) Service {
	return &service{
		repository:                    repository,
		investmentService:             investmentService,
		investmentDistributionService: investmentDistributionService,
		participantService:            participantService,
		benefitTypeService:            benefitTypeService,
		feeSettingService:             feeSettingService,
	}
}

func (s *service) Create(ctx context.Context, req *entity.BenefitParticipationDto) (*entity.BenefitParticipationDto, error) {
	req.CustomerID = shared.GetUserCredential(ctx).CustomerID
	req.OrganizationID = shared.GetOrganization(ctx).ID
	req.Status = model.BenefitParticipationStatusPending
	// Validate participant exists
	_, err := s.participantService.FindByID(ctx, req.ParticipantID)
	if err != nil {
		return nil, errors.New("participant not found")
	}

	// Validate benefit types in details
	for _, detail := range req.Details {
		_, err := s.benefitTypeService.FindByID(ctx, detail.BenefitTypeID)
		if err != nil {
			return nil, errors.New("benefit type not found")
		}
		detail.Status = model.BenefitParticipationStatusPending
	}

	existing, err := s.repository.FindByParticipantID(ctx, req.ParticipantID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return s.repository.Create(ctx, req, s.createInvestmentCallback)
	}

	return s.updateExistingBenefitParticipation(ctx, req, existing)
}

func (s *service) updateExistingBenefitParticipation(ctx context.Context, req *entity.BenefitParticipationDto, existing *entity.BenefitParticipationDto) (*entity.BenefitParticipationDto, error) {
	newDetails, pendingDetailsToReplace := s.categorizeDetails(req.Details, existing.Details)

	if len(newDetails) == 0 && len(pendingDetailsToReplace) == 0 {
		return existing, nil
	}

	err := s.repository.AppendOrUpdateDetails(ctx, existing.ID, newDetails, pendingDetailsToReplace, s.createInvestmentCallback)
	if err != nil {
		return nil, err
	}

	// Update existing details in memory
	for _, replacedDetail := range pendingDetailsToReplace {
		for i, existingDetail := range existing.Details {
			if existingDetail.ID == replacedDetail.ID {
				existing.Details[i] = replacedDetail
				break
			}
		}
	}
	existing.Details = append(existing.Details, newDetails...)

	return existing, nil
}

func (s *service) categorizeDetails(newDetails, existingDetails []*entity.BenefitParticipationDetailDto) ([]*entity.BenefitParticipationDetailDto, []*entity.BenefitParticipationDetailDto) {
	toAdd := make([]*entity.BenefitParticipationDetailDto, 0)
	toReplace := make([]*entity.BenefitParticipationDetailDto, 0)

	for _, newDetail := range newDetails {
		found := false
		for _, existingDetail := range existingDetails {
			if existingDetail.BenefitTypeID == newDetail.BenefitTypeID {
				if existingDetail.Status == model.BenefitParticipationStatusPending {
					newDetail.ID = existingDetail.ID
					toReplace = append(toReplace, newDetail)
				}
				found = true
				break
			}
		}
		if !found {
			toAdd = append(toAdd, newDetail)
		}
	}

	return toAdd, toReplace
}

func (s *service) createInvestmentCallback(ctx context.Context, tx *gorm.DB, dto *entity.BenefitParticipationDto) error {
	for _, detail := range dto.Details {
		if detail.MonthlyContribution > 0 {
			if err := s.createInvestmentForDetailWithTx(ctx, tx, dto, detail); err != nil {
				return err
			}
		}
	}

	if err := s.repository.UpdateWithTx(ctx, dto, tx); err != nil {
		return err
	}

	return nil
}

func (s *service) generateUniqueInvestmentCode(ctx context.Context) (string, error) {
	const maxRetries = 10

	for range maxRetries {
		// Generate 7 digit number (0000000 - 9999999)
		code, err := gonanoid.Generate("0123456789", 7)
		if err != nil {
			return "", err
		}

		_, err = s.investmentService.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
	}

	return "", errors.New("failed to generate unique investment code after maximum retries")
}

func (s *service) createInvestmentForDetailWithTx(ctx context.Context, tx *gorm.DB, benefitParticipation *entity.BenefitParticipationDto, detail *entity.BenefitParticipationDetailDto) error {
	distributions, err := s.investmentDistributionService.FindByParticipantID(ctx, benefitParticipation.ParticipantID)
	if err != nil {
		return err
	}

	if len(distributions) == 0 {
		return nil
	}

	investmentDto := &investmentEntity.InvestmentDto{
		OrganizationID:  benefitParticipation.OrganizationID,
		ParticipantID:   benefitParticipation.ParticipantID,
		CustomerID:      benefitParticipation.CustomerID,
		Amount:          detail.MonthlyContribution,
		Type:            model.InvestmentTypePPIP,
		InvestmentAt:    time.Now().Local().UTC(),
		Status:          model.InvestmentStatusCreated,
		InvestmentItems: make([]*investmentEntity.InvestmentItemDto, 0),
		Source:          model.InvestmentSourceBenefitParticipation,
	}

	investmentCode, err := s.generateUniqueInvestmentCode(ctx)
	if err != nil {
		return err
	}
	investmentDto.Code = investmentCode

	feeSetting, err := s.feeSettingService.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Create investment items based on distributions
	for _, distribution := range distributions {
		amount := detail.MonthlyContribution * distribution.Percent / 100
		feeAmount := amount * float64(feeSetting.AdminFee) / 100
		totalAmount := amount - feeAmount
		if amount > 0 {
			investmentItem := &investmentEntity.InvestmentItemDto{
				Amount:              amount,
				OrganizationID:      benefitParticipation.OrganizationID,
				ParticipantID:       benefitParticipation.ParticipantID,
				CustomerID:          benefitParticipation.CustomerID,
				InvestmentType:      model.InvestmentTypePPIP,
				InvestmentProductID: distribution.InvestmentProductID,
				Type:                model.InvestmentTypePPIP,
				FeeAmount:           feeAmount,
				TotalAmount:         totalAmount,
				Percent:             distribution.Percent,
				InvestmentAt:        time.Now().Local().UTC(),
				Status:              model.InvestmentStatusCreated,
			}
			investmentDto.InvestmentItems = append(investmentDto.InvestmentItems, investmentItem)
		}
	}

	investmentModel := investmentDto.ToModel()

	if err := tx.WithContext(ctx).Create(investmentModel).Error; err != nil {
		return err
	}

	benefitParticipation.InvestmentID = investmentModel.ID

	return nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.BenefitParticipationDto, error) {
	benefitParticipation, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return benefitParticipation, nil
}

func (s *service) FindByInvestmentPaymentID(ctx context.Context, investmentPaymentID string) (*entity.BenefitParticipationDto, error) {
	return s.repository.FindByInvestmentPaymentID(ctx, investmentPaymentID)
}

func (s *service) ConfirmationApprovalCallback(ctx context.Context, benefitParticipationID string, status model.BenefitParticipationStatus, tx *gorm.DB) error {
	err := s.repository.UpdateStatus(ctx, benefitParticipationID, status, tx)
	if err != nil {
		return err
	}

	err = s.repository.UpdateDetailsStatusByBenefitParticipationID(ctx, benefitParticipationID, status, tx)
	if err != nil {
		return err
	}

	return nil
}
