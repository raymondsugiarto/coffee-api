package investment

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	ec "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	feesetting "github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	netassetvalue "github.com/raymondsugiarto/coffee-api/pkg/module/net_asset_value"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const ONE_ITEM = 1

type Service interface {
	Create(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
	FindByID(ctx context.Context, id string) (*entity.InvestmentDto, error)
	FindByCode(ctx context.Context, code string) (*entity.InvestmentDto, error)
	Update(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error)
	ConfirmationApprovalCallback(ctx context.Context, req *entity.InvestmentPaymentDto, tx *gorm.DB) (context.Context, error)
	NotifyApprovalCallback(ctx context.Context, req *entity.InvestmentPaymentDto) error
	CreateFromPaymentFee(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error)
	UploadPayment(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
}

type service struct {
	repo                          Repository
	participantService            participant.Service
	investmentItemService         investmentitem.Service
	investmentPaymentService      investmentpayment.Service
	investmentDistributionService investmentdistribution.Service
	unitLinkService               unitlink.Service
	netAssetValueService          netassetvalue.Service
	notificationService           notification.Service
	approvalService               approval.Service
	feeSettingService             feesetting.Service
	customerService               customer.Service
}

func NewService(
	repo Repository,
	participantService participant.Service,
	investmentItemService investmentitem.Service,
	investmentPaymentService investmentpayment.Service,
	investmentDistributionService investmentdistribution.Service,
	unitLinkService unitlink.Service,
	netAssetValueService netassetvalue.Service,
	notificationService notification.Service,
	approvalService approval.Service,
	feeSettingService feesetting.Service,
	customerService customer.Service,
) Service {
	return &service{
		repo:                          repo,
		participantService:            participantService,
		investmentItemService:         investmentItemService,
		investmentPaymentService:      investmentPaymentService,
		investmentDistributionService: investmentDistributionService,
		unitLinkService:               unitLinkService,
		netAssetValueService:          netAssetValueService,
		notificationService:           notificationService,
		approvalService:               approvalService,
		feeSettingService:             feeSettingService,
		customerService:               customerService,
	}
}

func calculateInvestmentAmounts(baseAmount, percent float64, adminFee float64) (grossAmount, feeAmount, netInvestmentAmount float64) {
	netInvestmentAmount = baseAmount * percent / 100
	feeAmount = netInvestmentAmount * adminFee / 100
	// Gross amount is the total payment required (net investment + fee)
	grossAmount = netInvestmentAmount + feeAmount

	return
}

func (s *service) Create(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID

	dto.Type = model.InvestmentTypePPIP
	if dto.CompanyID != "" {
		dto.Type = model.InvestmentTypeDKP

		// validate if company id has already set investment distribution
		if err := s.investmentDistributionService.HaveDistribution(ctx, dto.CompanyID); err != nil {
			return nil, status.New(status.BadRequest, fmt.Errorf("perusahaan belum mempunyai alokasi investasi"))
		}
	} else {
		dto.InvestmentAt = time.Now().Local().UTC()
	}

	// companyID := shared.GetCompanyID(ctx)

	dto.ExpiredAt = time.Now().Add(time.Minute * 60).Local().UTC()
	if dto.Type == model.InvestmentTypeDKP {
		return s.createDKP(ctx, dto)
	} else {
		return s.createPPIP(ctx, dto)
	}
}

func (s *service) CreateFromPaymentFee(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	dto.CustomerID = shared.GetUserCredential(ctx).CustomerID
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.Type = model.InvestmentTypePPIP
	dto.InvestmentAt = time.Now().Local().UTC()
	dto.Status = model.InvestmentStatusSuccess
	dto.ExpiredAt = time.Now().Add(time.Minute * 60).Local().UTC()

	investmentProducts, err := s.unitLinkService.FindAllInvestmentProductByParticipant(ctx, dto.ParticipantID)
	if err != nil {
		return nil, err
	}
	dto.Status = model.InvestmentStatusCreated

	feeSetting, err := s.feeSettingService.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	investmentItems := make([]*entity.InvestmentItemDto, 0)
	for _, item := range investmentProducts {
		// Calculate fees properly based on the existing TotalAmount (treat it as net investment amount)
		// Since item.TotalAmount is the existing investment amount, we need to calculate the gross amount with fees
		netAmount := item.TotalAmount
		feeAmount := netAmount * feeSetting.AdminFee / 100
		grossAmount := netAmount + feeAmount

		investmentItem := &entity.InvestmentItemDto{
			EmployeeAmount:      grossAmount,
			Amount:              grossAmount, // Gross amount (net + fee)
			OrganizationID:      dto.OrganizationID,
			ParticipantID:       dto.ParticipantID,
			CustomerID:          dto.CustomerID,
			InvestmentType:      dto.Type,
			InvestmentProductID: item.InvestmentProductID,
			Type:                dto.Type,
			FeeAmount:           feeAmount, // Calculated fee amount
			TotalAmount:         netAmount, // Net investment amount
			Percent:             -1,        // -1 means not applicable
			ExpiredAt:           dto.ExpiredAt,
			InvestmentAt:        dto.InvestmentAt,
			Status:              dto.Status,
		}
		investmentItems = append(investmentItems, investmentItem)
	}

	dto.InvestmentItems = investmentItems

	return s.create(ctx, dto)
}

func (s *service) UploadPayment(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	investment, err := s.repo.Get(ctx, dto.InvestmentID)
	if err != nil {
		return nil, err
	}

	dto.Amount = math.Ceil(investment.GetTotalPaymentAmount()) // the total payment amount is rounded up to avoid decimals in the user's payment
	dto.Status = model.InvestmentPaymentStatusPending
	dto, err = s.investmentPaymentService.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.UpdateStatusWaitingVerification(ctx, &entity.InvestmentDto{
		ID: dto.InvestmentID,
	})
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (s *service) createDKP(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	// TODO: get all employees from company

	feeSetting, err := s.feeSettingService.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	participants, err := s.participantService.FindAllParticipantCompany(ctx, &e.ParticipantFindAllRequest{
		InvestmentAt: dto.InvestmentAt,
	})
	if err != nil {
		return dto, err
	}

	// get distribution by company
	distributions, err := s.investmentDistributionService.FindByCompanyID(ctx, dto.CompanyID)
	if err != nil {
		return dto, err
	}

	totalExpectedAmount := 0.0
	for _, participant := range participants {
		totalExpectedAmount += *participant.Customer.EmployerAmount +
			*participant.Customer.EmployeeAmount +
			*participant.Customer.VoluntaryAmount +
			*participant.Customer.EducationFundAmount
	}

	// Step 2: Hitung selisih rasio berdasarkan input manual dari user
	diffRatio := 0.0
	if totalExpectedAmount > 0 {
		diffRatio = dto.Amount / totalExpectedAmount
	}

	// loop all employees, get participant = DKP and create investment (per participant)
	investmentItems := make([]*entity.InvestmentItemDto, 0)
	for _, participant := range participants {
		var baseAmount float64
		var employerAmount, employeeAmount, voluntaryAmount, educationFundAmount float64

		if participant.Customer.EmployerAmount != nil {
			employerAmount = *participant.Customer.EmployerAmount
			baseAmount += employerAmount
		}
		if participant.Customer.EmployeeAmount != nil {
			employeeAmount = *participant.Customer.EmployeeAmount
			baseAmount += employeeAmount
		}
		if participant.Customer.VoluntaryAmount != nil {
			voluntaryAmount = *participant.Customer.VoluntaryAmount
			baseAmount += voluntaryAmount
		}
		if participant.Customer.EducationFundAmount != nil {
			educationFundAmount = *participant.Customer.EducationFundAmount
			baseAmount += educationFundAmount
		}
		adjustedBaseAmount := baseAmount * diffRatio

		// Apply adjustment ratio to individual amounts proportionally
		adjustmentFactor := diffRatio
		adjustedEmployerAmount := employerAmount * adjustmentFactor
		adjustedEmployeeAmount := employeeAmount * adjustmentFactor
		adjustedVoluntaryAmount := voluntaryAmount * adjustmentFactor
		adjustedEducationFundAmount := educationFundAmount * adjustmentFactor
		for _, distribution := range distributions {
			if distribution.Percent <= 0 {
				continue // Skip if percent is 0
			}
			if len(distributions) == ONE_ITEM {
				distribution.Percent = 100
			}

			grossAmount, feeAmount, netInvestmentAmount := calculateInvestmentAmounts(adjustedBaseAmount, distribution.Percent, feeSetting.AdminFee)

			// Calculate proportional amounts for this investment item based on percentage
			itemEmployerAmount := adjustedEmployerAmount * distribution.Percent / 100
			itemEmployeeAmount := adjustedEmployeeAmount * distribution.Percent / 100
			itemVoluntaryAmount := adjustedVoluntaryAmount * distribution.Percent / 100
			itemEducationFundAmount := adjustedEducationFundAmount * distribution.Percent / 100

			investmentItem := &entity.InvestmentItemDto{
				Amount:              grossAmount,
				OrganizationID:      dto.OrganizationID,
				ParticipantID:       participant.ID,
				CustomerID:          participant.CustomerID,
				InvestmentType:      dto.Type,
				InvestmentProductID: distribution.InvestmentProductID,
				Type:                dto.Type,
				Percent:             distribution.Percent,
				FeeAmount:           feeAmount,
				TotalAmount:         netInvestmentAmount,
				EmployerAmount:      itemEmployerAmount,
				EmployeeAmount:      itemEmployeeAmount,
				VoluntaryAmount:     itemVoluntaryAmount,
				EducationFundAmount: itemEducationFundAmount,
				ExpiredAt:           dto.ExpiredAt,
				InvestmentAt:        dto.InvestmentAt,
				Status:              dto.Status,
			}
			investmentItems = append(investmentItems, investmentItem)
		}
	}
	dto.InvestmentItems = investmentItems

	// if investment from company, then automatically create investment payment and status request
	dto.Status = model.InvestmentStatusRequest

	return s.create(ctx, dto)
}

func (s *service) createPPIP(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	// hanya dari mobile
	if dto.IsNewParticipant {
		participantDto := new(e.ParticipantDto)
		// get customer id from context
		customerID := shared.GetUserCredential(ctx).CustomerID
		participantDto.CustomerID = customerID
		participantDto.OrganizationID = shared.GetOrganization(ctx).ID

		// Get customer information to determine participant type
		customer, err := s.customerService.FindByID(ctx, customerID)
		if err != nil {
			return dto, err
		}

		participant, err := s.participantService.Create(ctx, participantDto, customer)
		if err != nil {
			return dto, err
		}
		dto.CustomerID = participant.CustomerID
		dto.ParticipantID = participant.ID

		// fill participantID to investment items if new participant
		for _, item := range dto.InvestmentItems {
			item.ParticipantID = participant.ID
			item.CustomerID = participant.CustomerID
		}

		// create investment distribution
		bgCtx := shared.NewBackgroundContext(ctx)
		go s.investmentDistributionService.CreateBatch(bgCtx, dto.ToInvestmentDistributions())
	}
	dto.Status = model.InvestmentStatusCreated

	feeSetting, err := s.feeSettingService.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	investmentItems := make([]*entity.InvestmentItemDto, 0)
	for _, item := range dto.InvestmentItems {

		if len(dto.InvestmentItems) == ONE_ITEM {
			item.Percent = 100
		}
		grossAmount, feeAmount, netInvestmentAmount := calculateInvestmentAmounts(dto.Amount, item.Percent, feeSetting.AdminFee)

		investmentItem := &entity.InvestmentItemDto{
			EmployeeAmount:      grossAmount,
			Amount:              grossAmount,
			OrganizationID:      dto.OrganizationID,
			ParticipantID:       dto.ParticipantID,
			CustomerID:          dto.CustomerID,
			InvestmentType:      dto.Type,
			InvestmentProductID: item.InvestmentProductID,
			Type:                dto.Type,
			FeeAmount:           feeAmount,
			TotalAmount:         netInvestmentAmount,
			Percent:             item.Percent,
			ExpiredAt:           dto.ExpiredAt,
			InvestmentAt:        dto.InvestmentAt,
			Status:              dto.Status,
		}
		investmentItems = append(investmentItems, investmentItem)
	}

	dto.InvestmentItems = investmentItems

	return s.create(ctx, dto)
}

func (s *service) generateUniqueInvestmentCode(ctx context.Context) (string, error) {
	const maxRetries = 10

	for range maxRetries {
		// Generate 7 digit number (0000000 - 9999999)
		code, err := gonanoid.Generate("0123456789", 7)
		if err != nil {
			return "", err
		}

		_, err = s.repo.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return code, nil
			}
			return "", err
		}
	}

	return "", errors.New("failed to generate unique investment code after maximum retries")
}

func (s *service) FindByCode(ctx context.Context, code string) (*entity.InvestmentDto, error) {
	return s.repo.FindByCode(ctx, code)
}

func (s *service) create(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	investmentCode, err := s.generateUniqueInvestmentCode(ctx)
	if err != nil {
		return nil, err
	}
	dto.Code = investmentCode

	// clear investment items, insert separate
	dtoCreate := *dto
	dtoCreate.InvestmentItems = make([]*entity.InvestmentItemDto, 0)
	res, err := s.repo.Create(ctx, &dtoCreate, func(ctx context.Context, db *gorm.DB, dtoCreate *entity.InvestmentDto) error {

		// fill investment id for each investment item
		investmentItems := lo.Map(dto.InvestmentItems, func(item *entity.InvestmentItemDto, _ int) *entity.InvestmentItemDto {
			item.InvestmentID = &dtoCreate.ID
			return item
		})

		// create investment items
		_, err := s.investmentItemService.CreateBatchWithTx(ctx, db, investmentItems)
		if err != nil {
			return err
		}

		// set investment items for calculate total payment amount
		dtoCreate.InvestmentItems = investmentItems
		if dtoCreate.Type == model.InvestmentTypeDKP {
			paymentDto := dtoCreate.ToCompanyInvestmentPayment()
			if confirmationImage, ok := ctx.Value("confirmationImage").(string); ok {
				paymentDto.ConfirmationImageUrl = confirmationImage
			}
			paymentCreated, err := s.investmentPaymentService.Create(ctx, paymentDto)
			if err != nil {
				return err
			}
			uid := shared.GetUserCredential(ctx).UserID
			paymentDto.Investment = dtoCreate
			paymentDto.ID = paymentCreated.ID
			if _, err := s.approvalService.Create(ctx, paymentDto.ToApprovalSubmitDto(uid, model.ApprovalTypeInvestment)); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	bgCtx := shared.NewBackgroundContext(ctx)
	// go s.customerService.UpdateEffectiveDateIfFirstPayment(bgCtx, res.ID, res.InvestmentAt)
	go s.notificationService.NotifyInvestment(bgCtx, res)

	return res, err
}

type ctxKey string

const unitLinkDtosKey ctxKey = "unitLinkDtos"

func (s *service) NotifyApprovalCallback(ctx context.Context, req *entity.InvestmentPaymentDto) error {
	// get unit link dtos from context
	unitLinkDtos, ok := ctx.Value(unitLinkDtosKey).([]*ec.UnitLinkDto)
	if !ok || len(unitLinkDtos) == 0 {
		log.WithContext(ctx).Error("No unit link dtos found in context")
		return nil
	}

	bgCtx := shared.NewBackgroundContext(ctx)
	go s.netAssetValueService.PublishByInvestment(bgCtx, unitLinkDtos)

	return nil
}
func (s *service) ConfirmationApprovalCallback(ctx context.Context, req *entity.InvestmentPaymentDto, tx *gorm.DB) (context.Context, error) {
	investmentPayment, err := s.investmentPaymentService.FindByID(ctx, req.ID)
	if err != nil {
		return ctx, err
	}

	investmentPayment.Status = req.Status
	s.investmentPaymentService.Update(ctx, investmentPayment)

	if investmentPayment.Status != model.InvestmentPaymentStatusSuccess {
		return ctx, nil
	}

	// create unit link
	investment, err := s.repo.Get(ctx, investmentPayment.InvestmentID)
	if err != nil {
		return ctx, err
	}

	investment.Status = model.InvestmentStatusSuccess
	_, err = s.repo.UpdateWithTx(ctx, tx, investment)
	if err != nil {
		return ctx, err
	}

	participantIDs := make(map[string]bool)
	for _, item := range investment.InvestmentItems {
		if item.ParticipantID != "" {
			participantIDs[item.ParticipantID] = true
		}
	}

	for participantID := range participantIDs {
		err = s.participantService.UpdateStatus(ctx, participantID, model.ParticipantStatusActive)
		if err != nil {
			log.WithContext(ctx).Error("Failed to update participant status", "participantID", participantID, "error", err)
		}
	}

	log.WithContext(ctx).Info("Processing investment approval", "id", investment.ID, "type", investment.Type)

	// Add approval date to be used as transaction date for UnitLink records
	approvalDate := time.Now().Local().UTC()

	unitLinkDtos := make([]*ec.UnitLinkDto, 0)
	for _, item := range investment.InvestmentItems {
		if item.ParticipantID == "" {
			log.WithContext(ctx).Warn("Investment item missing participant ID", "itemID", item.ID)
			continue
		}

		unitLinkDto := &ec.UnitLinkDto{
			OrganizationID:      item.OrganizationID,
			CustomerID:          item.CustomerID,
			ParticipantID:       item.ParticipantID,
			InvestmentProductID: item.InvestmentProductID,
			Type:                item.Type,
			TotalAmount:         item.TotalAmount, // Use TotalAmount (after fee deduction)
			TransactionDate:     approvalDate,
		}
		unitLinkDtos = append(unitLinkDtos, unitLinkDto)
	}

	for _, item := range unitLinkDtos {
		log.WithContext(ctx).Info("Creating unit link", "participantID", item.ParticipantID)

		res, err := s.unitLinkService.CreateWithTx(ctx, tx, item)
		if err != nil {
			log.WithContext(ctx).Error("Failed to create unit link", "participantID", item.ParticipantID, "error", err)
			return ctx, err
		}
		item.TotalAmount = res.TotalAmount
		item.ID = res.ID
	}

	// save context with unit link dto
	ctx = context.WithValue(ctx, unitLinkDtosKey, unitLinkDtos)
	bgCtx := shared.NewBackgroundContext(ctx)
	go s.customerService.UpdateEffectiveDateIfFirstPayment(bgCtx, investment.ID, time.Now())

	return ctx, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.InvestmentDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.InvestmentDto) (*entity.InvestmentDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}
