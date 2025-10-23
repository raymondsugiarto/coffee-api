package entity

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type InvestmentPaymentConfirmationInputDto struct {
	InvestmentPaymentID string
	PaymentAt           time.Time
	ConfirmationImage   *multipart.FileHeader `json:"confirmationImage"`
}

func (d *InvestmentPaymentConfirmationInputDto) ToDto() *InvestmentPaymentDto {
	dto := &InvestmentPaymentDto{
		ID:        d.InvestmentPaymentID,
		PaymentAt: d.PaymentAt,
	}

	if d.ConfirmationImage != nil {
		dto.ConfirmationImageUrl = d.ConfirmationImage.Filename
	}

	return dto
}

type InvestmentPaymentInputDto struct {
	InvestmentID string
	BankID       string
}

func (d *InvestmentPaymentInputDto) ToDto() *InvestmentPaymentDto {
	return &InvestmentPaymentDto{
		InvestmentID: d.InvestmentID,
		BankID:       d.BankID,
	}
}

type InvestmentPaymentDto struct {
	ID                   string                        `json:"id"`
	OrganizationID       string                        `json:"organizationId"`
	InvestmentID         string                        `json:"investmentId"`
	Investment           *InvestmentDto                `json:"investment"`
	BankID               string                        `json:"bankId"`
	Bank                 *entity.BankDto               `json:"bank"`
	BankCode             string                        `json:"bankCode"`
	BankName             string                        `json:"bankName"`
	AccountName          string                        `json:"accountName"`
	AccountNumber        string                        `json:"accountNumber"`
	Amount               float64                       `json:"amount"`
	ConfirmationImageUrl string                        `json:"confirmationImageUrl"`
	PaymentAt            time.Time                     `json:"paymentAt"`
	Status               model.InvestmentPaymentStatus `json:"status"`
	CreatedAt            time.Time                     `json:"createdAt"`
	UpdatedAt            time.Time                     `json:"updatedAt"`
}

type UnitLinkDto struct {
	ID                  string                 `json:"id"`
	OrganizationID      string                 `json:"-"`
	TransactionDate     time.Time              `json:"transactionDate"`
	ParticipantID       string                 `json:"participantId"`
	Participant         *entity.ParticipantDto `json:"participant"`
	CustomerID          string                 `json:"customerId"`
	Customer            *entity.CustomerDto    `json:"customer"`
	InvestmentProductID string                 `json:"investmentProductId"`
	InvestmentProduct   *InvestmentProductDto  `json:"investmentProduct"`
	Type                model.InvestmentType   `json:"type"`
	TotalAmount         float64                `json:"totalAmount"`
	Nab                 float64                `json:"nab"`
	Ip                  float64                `json:"ip"`
	CreatedAt           time.Time              `json:"createdAt"`
}

func (d *InvestmentPaymentDto) ToModel() *model.InvestmentPayment {
	m := &model.InvestmentPayment{
		OrganizationID:       d.OrganizationID,
		InvestmentID:         d.InvestmentID,
		BankID:               d.BankID,
		BankCode:             d.BankCode,
		BankName:             d.BankName,
		AccountName:          d.AccountName,
		AccountNumber:        d.AccountNumber,
		Amount:               d.Amount,
		ConfirmationImageUrl: d.ConfirmationImageUrl,
		PaymentAt:            d.PaymentAt,
		Status:               d.Status,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *InvestmentPaymentDto) PaymentConfirmationToModel() *model.InvestmentPayment {
	m := &model.InvestmentPayment{
		PaymentAt:            d.PaymentAt,
		Status:               d.Status,
		ConfirmationImageUrl: d.ConfirmationImageUrl,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *InvestmentPaymentDto) FromModel(m *model.InvestmentPayment) *InvestmentPaymentDto {
	p := &InvestmentPaymentDto{
		ID:                   m.CommonWithIDs.ID,
		OrganizationID:       m.OrganizationID,
		InvestmentID:         m.InvestmentID,
		BankID:               m.BankID,
		BankCode:             m.BankCode,
		BankName:             m.BankName,
		AccountName:          m.AccountName,
		AccountNumber:        m.AccountNumber,
		Amount:               m.Amount,
		ConfirmationImageUrl: m.ConfirmationImageUrl,
		PaymentAt:            m.PaymentAt,
		Status:               m.Status,
		CreatedAt:            m.CommonWithIDs.CreatedAt,
		UpdatedAt:            m.CommonWithIDs.UpdatedAt,
	}

	if m.Investment != nil {
		p.Investment = new(InvestmentDto).FromModel(m.Investment)
	}

	return p
}

func (d *InvestmentPaymentDto) ToApprovalSubmitDto(uid string, approvalType model.ApprovalType) *entity.ApprovalDto {
	return &entity.ApprovalDto{
		OrganizationID: d.OrganizationID,
		UserIDRequest:  uid,
		RefID:          d.ID,
		RefTable:       "investment_payment",
		Detail:         fmt.Sprintf("Pembayaran Investasi [%s]", d.Investment.Code),
		Type:           approvalType,
		Action:         "PAYMENT",
		Status:         "SUBMIT",
		Reason:         "New Investment Payment",
	}
}

type InvestmentPaymentCompanyGetMonthlyDto struct {
	TotalMonthlyCompanyContribution float64 `json:"totalMonthlyCompanyContribution"`
	TotalEmployeeAmount             float64 `json:"totalEmployeeAmount"`
	TotalEmployerAmount             float64 `json:"totalEmployerAmount"`
	TotalVoluntaryAmount            float64 `json:"totalVoluntaryAmount"`
	TotalEducationFundAmount        float64 `json:"totalEducationFundAmount"`
	CountEmployeeUnpaid             int     `json:"countEmployeeUnpaid"`
	CountEmployeePaid               int     `json:"countEmployeePaid"`
	Fee                             float64 `json:"fee"`
}

type InvestmentPaymentFindAllRequest struct {
	entity.FindAllRequest
	CustomerID           string
	CompanyID            string
	ShowAll              bool
	StartDate            time.Time
	EndDate              time.Time
	IncludeNetAssetValue bool
}

func (r *InvestmentPaymentFindAllRequest) GenerateFilter() {
	if !r.ShowAll {
		if !r.StartDate.IsZero() {
			r.FindAllRequest.AddFilter(pagination.FilterItem{
				Field: "payment_at",
				Op:    "gte",
				Val:   r.StartDate,
			})
		}
		if !r.EndDate.IsZero() {
			r.FindAllRequest.AddFilter(pagination.FilterItem{
				Field: "payment_at",
				Op:    "lt",
				Val:   r.EndDate,
			})
		}
	}
}

type InvestmentPaymentSummaryDto struct {
	TotalPayments         float64 `json:"totalPayments"`
	TotalApprovedPayments float64 `json:"totalApprovedPayments"`
	TotalPendingPayments  float64 `json:"totalPendingPayments"`
}

type GetTotalMonthlyCompanyContributionRequest struct {
	entity.FindAllRequest
	InvestmentAt               time.Time
	CalculateAll               bool
	IncludeEmployeeInformation bool
	UsePeriod                  bool // if true, use the month and year of InvestmentAt, otherwise use the full date
}

func (i *InvestmentPaymentDto) GetInfo() entity.RejectEmail {
	return entity.RejectEmail{
		Email:       i.Investment.Company.Email,
		Name:        i.Investment.Code,
		Description: "Pengajuan Investasi Dengan Kode " + i.Investment.Code,
	}
}
