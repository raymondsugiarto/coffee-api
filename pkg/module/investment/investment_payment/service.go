package investmentpayment

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	"github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
	FindByID(ctx context.Context, id string) (*entity.InvestmentPaymentDto, error)
	Update(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
	UpdatePaymentConfirmation(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest) (*pagination.ResultPagination, error)
	GenerateExcel(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest, result *pagination.ResultPagination) ([]byte, error)
	GetTotalMonthlyCompanyContribution(ctx context.Context, req *entity.GetTotalMonthlyCompanyContributionRequest) (*entity.InvestmentPaymentCompanyGetMonthlyDto, error)
	GetPaymentSummary(ctx context.Context) (*entity.InvestmentPaymentSummaryDto, error)
}

type service struct {
	repo                Repository
	customerService     customer.Service
	notificationService notification.Service
	participantService  participant.Service
	approvalService     approval.Service
	feeSetting          fee_setting.Service
}

func NewService(repo Repository,
	customerService customer.Service,
	notificationService notification.Service,
	participantService participant.Service,
	approvalService approval.Service,
	feeSetting fee_setting.Service,
) Service {
	return &service{repo,
		customerService,
		notificationService,
		participantService,
		approvalService,
		feeSetting,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.InvestmentPaymentDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) UpdatePaymentConfirmation(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	dto.Status = model.InvestmentPaymentStatusConfirmed
	dto.PaymentAt = time.Now()

	payment, err := s.repo.Get(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	uid := shared.GetUserCredential(ctx).UserID

	approvalType := model.ApprovalTypeInvestment
	if payment.Investment != nil && payment.Investment.Source == model.InvestmentSourceBenefitParticipation {
		approvalType = model.ApprovalTypeBenefitParticipation
	}

	result, err := s.repo.UpdatePaymentConfirmation(ctx, dto, func(tx *gorm.DB) error {
		if _, err := s.approvalService.CreateWithTx(ctx, payment.ToApprovalSubmitDto(uid, approvalType), tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	bgCtx := shared.NewBackgroundContext(ctx)
	go s.notificationService.NotifyInvestmentPaymentConfirmed(bgCtx, payment.Investment)

	return result, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) GetTotalMonthlyCompanyContribution(ctx context.Context, req *entity.GetTotalMonthlyCompanyContributionRequest) (*entity.InvestmentPaymentCompanyGetMonthlyDto, error) {
	// Call the repository method to get the total monthly company contribution
	result := new(entity.InvestmentPaymentCompanyGetMonthlyDto)

	if req.IncludeEmployeeInformation {
		paidParticipants, err := s.participantService.FindAllParticipantCompany(ctx, &e.ParticipantFindAllRequest{
			InvestmentAt: req.InvestmentAt,
			CalculateAll: req.CalculateAll,
			PaidEmployee: true,
			UsePeriod:    req.UsePeriod,
		})
		if err != nil {
			return nil, err
		}
		result.CountEmployeePaid = len(paidParticipants)

		fee, err := s.feeSetting.GetConfig(ctx)
		if err != nil {
			return nil, err
		}

		if fee != nil {
			result.Fee = fee.AdminFee
		} else {
			result.Fee = 0
		}
	}
	participants, err := s.participantService.FindAllParticipantCompany(ctx, &e.ParticipantFindAllRequest{
		InvestmentAt: req.InvestmentAt,
		CalculateAll: req.CalculateAll,
		PaidEmployee: false,
		UsePeriod:    req.UsePeriod,
	})
	if err != nil {
		return nil, err
	}
	result.CountEmployeeUnpaid = len(participants)
	var (
		totalEmployeeAmt      float64
		totalEmployerAmt      float64
		totalVoluntaryAmt     float64
		totalEducationFundAmt float64
	)
	for _, participant := range participants {
		if participant.Customer != nil {
			log.WithContext(ctx).Infof("Participant: %s, EmployerAmount: %d, EmployeeAmount: %d, VoluntaryAmount: %d, EducationFundAmount: %d", &participant.Code, &participant.Customer.EmployerAmount, &participant.Customer.EmployeeAmount, &participant.Customer.VoluntaryAmount, &participant.Customer.EducationFundAmount)
			totalEmployeeAmt += *participant.Customer.EmployeeAmount
			totalEmployerAmt += *participant.Customer.EmployerAmount
			totalVoluntaryAmt += *participant.Customer.VoluntaryAmount
			totalEducationFundAmt += *participant.Customer.EducationFundAmount
		}
	}

	result.TotalEmployeeAmount = totalEmployeeAmt
	result.TotalEmployerAmount = totalEmployerAmt
	result.TotalVoluntaryAmount = totalVoluntaryAmt
	result.TotalEducationFundAmount = totalEducationFundAmt
	result.TotalMonthlyCompanyContribution = totalEmployeeAmt + totalEmployerAmt + totalVoluntaryAmt + totalEducationFundAmt
	return result, nil
}

func getStatusInIndonesian(status model.InvestmentPaymentStatus) string {
	switch status {
	case model.InvestmentPaymentStatusPending:
		return "Menunggu"
	case model.InvestmentPaymentStatusConfirmed:
		return "Dikonfirmasi"
	case model.InvestmentPaymentStatusRejected:
		return "Ditolak"
	case model.InvestmentPaymentStatusSuccess:
		return "Berhasil"
	case model.InvestmentPaymentStatusExpired:
		return "Kedaluwarsa"
	default:
		return string(status)
	}
}

func (s *service) GenerateExcel(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest, result *pagination.ResultPagination) ([]byte, error) {
	excel := excelize.NewFile()
	sheetName := "Investment Payments"
	index, _ := excel.NewSheet(sheetName)
	excel.DeleteSheet("Sheet1")

	headers := []string{
		"ID Pembayaran", "Kode Investasi", "Kode Partisipan", "Kode Perusahaan", "Nama Perusahaan",
		"Tipe Investasi", "Jumlah", "Status Pembayaran", "Tanggal Pembayaran",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		excel.SetCellValue(sheetName, cell, header)
	}

	payments := result.Data.([]*entity.InvestmentPaymentDto)
	for i, payment := range payments {
		row := i + 2 // Start from row 2 to skip header

		col1, _ := excelize.CoordinatesToCellName(1, row)
		excel.SetCellValue(sheetName, col1, payment.ID)

		col2, _ := excelize.CoordinatesToCellName(2, row)
		if payment.Investment != nil {
			excel.SetCellValue(sheetName, col2, payment.Investment.Code)
		} else {
			excel.SetCellValue(sheetName, col2, "")
		}

		col3, _ := excelize.CoordinatesToCellName(3, row)
		if payment.Investment != nil && payment.Investment.Participant != nil {
			excel.SetCellValue(sheetName, col3, payment.Investment.Participant.Code)
		} else {
			excel.SetCellValue(sheetName, col3, "")
		}

		col4, _ := excelize.CoordinatesToCellName(4, row)
		if payment.Investment != nil && payment.Investment.Company != nil {
			excel.SetCellValue(sheetName, col4, payment.Investment.Company.CompanyCode)
		} else {
			excel.SetCellValue(sheetName, col4, "")
		}

		col5, _ := excelize.CoordinatesToCellName(5, row)
		if payment.Investment != nil && payment.Investment.Company != nil {
			excel.SetCellValue(sheetName, col5, payment.Investment.Company.FirstName)
		} else {
			excel.SetCellValue(sheetName, col5, "")
		}

		col6, _ := excelize.CoordinatesToCellName(6, row)
		if payment.Investment != nil {
			excel.SetCellValue(sheetName, col6, string(payment.Investment.Type))
		} else {
			excel.SetCellValue(sheetName, col6, "")
		}

		col7, _ := excelize.CoordinatesToCellName(7, row)
		excel.SetCellValue(sheetName, col7, payment.Amount)

		col8, _ := excelize.CoordinatesToCellName(8, row)
		excel.SetCellValue(sheetName, col8, getStatusInIndonesian(payment.Status))

		col9, _ := excelize.CoordinatesToCellName(9, row)
		if !payment.PaymentAt.IsZero() {
			excel.SetCellValue(sheetName, col9, payment.PaymentAt.Format("2006-01-02 15:04:05"))
		} else {
			excel.SetCellValue(sheetName, col9, "")
		}
	}

	// Set active sheet
	excel.SetActiveSheet(index)

	// Write to buffer
	buffer, err := excel.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (s *service) GetPaymentSummary(ctx context.Context) (*entity.InvestmentPaymentSummaryDto, error) {
	companyID := shared.GetCompanyID(ctx)
	return s.repo.GetPaymentSummary(ctx, companyID)
}
