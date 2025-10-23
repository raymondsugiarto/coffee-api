package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type SumUnitLinkPortfolioDto struct {
	TotalAmount      float64 `json:"totalAmount"`
	TotalModal       float64 `json:"totalModal"`
	TotalUnit        float64 `json:"totalUnit"`
	CurrentBalance   float64 `json:"currentBalance"`
	Profit           float64 `json:"profit"`
	ReturnPercentage float64 `json:"returnPercentage"`
	Count            float64 `json:"count"`
}

type UnitLinkPortfolioDto struct {
	InvestmentProductID string                 `json:"investmentProductId"`
	InvestmentProduct   e.InvestmentProductDto `json:"investmentProduct"`
	TotalAmount         float64                `json:"totalAmount"`
	Nab                 float64                `json:"nab"`
	Ip                  float64                `json:"ip"`
}

type UnitLinkPortfolioGroupParticipantDto struct {
	OrganizationID      string  `json:"organizationId"`
	InvestmentProductID string  `json:"investmentProductId"`
	ParticipantID       string  `json:"participantId"`
	Type                string  `json:"type"`
	TotalAmount         float64 `json:"totalAmount"`
	Nab                 float64 `json:"nab"`
	Ip                  float64 `json:"ip"`
}

type UnitLinkLatestEachProductAndParticipantAndTypeDto struct {
	ID                  string    `json:"id"`
	InvestmentProductID string    `json:"investmentProductId"`
	ParticipantID       string    `json:"participantId"`
	Type                string    `json:"type"`
	TransactionDate     time.Time `json:"transactionDate"`
	Ip                  float64   `json:"ip"`
}

type UnitLinkSummaryCompanyDto struct {
	SumIp       float64 `json:"sumIp"`
	TotalAmount float64 `json:"totalAmount"`
}

type UnitLinkDto struct {
	ID                  string                  `json:"id"`
	OrganizationID      string                  `json:"-"`
	TransactionDate     time.Time               `json:"transactionDate"`
	ParticipantID       string                  `json:"participantId"`
	Participant         *entity.ParticipantDto  `json:"participant"`
	CustomerID          string                  `json:"customerId"`
	Customer            *entity.CustomerDto     `json:"customer"`
	InvestmentProductID string                  `json:"investmentProductId"`
	InvestmentProduct   *e.InvestmentProductDto `json:"investmentProduct"`
	Type                model.InvestmentType    `json:"type"`
	TotalAmount         float64                 `json:"totalAmount"`
	Nab                 float64                 `json:"nab"`
	Ip                  float64                 `json:"ip"`
	CreatedAt           time.Time               `json:"createdAt"`
}

func (u *UnitLinkDto) ToModel() *model.UnitLink {
	m := &model.UnitLink{
		ParticipantID:       u.ParticipantID,
		TransactionDate:     u.TransactionDate,
		OrganizationID:      u.OrganizationID,
		CustomerID:          u.CustomerID,
		InvestmentProductID: u.InvestmentProductID,
		Type:                u.Type,
		TotalAmount:         u.TotalAmount,
		Nab:                 u.Nab,
		Ip:                  u.Ip,
	}
	if u.ID != "" {
		m.ID = u.ID
	}
	return m
}

func (u *UnitLinkDto) FromModel(m *model.UnitLink) *UnitLinkDto {
	if m == nil {
		return nil
	}
	dto := &UnitLinkDto{
		ID:                  m.ID,
		ParticipantID:       m.ParticipantID,
		TransactionDate:     m.TransactionDate,
		OrganizationID:      m.OrganizationID,
		CustomerID:          m.CustomerID,
		InvestmentProductID: m.InvestmentProductID,
		Type:                m.Type,
		TotalAmount:         m.TotalAmount,
		Nab:                 m.Nab,
		Ip:                  m.Ip,
		CreatedAt:           m.CreatedAt,
	}
	if m.Customer != nil {
		dto.Customer = new(entity.CustomerDto).FromModel(m.Customer)
	}
	if m.Participant != nil {
		dto.Participant = new(entity.ParticipantDto).FromModel(m.Participant)
	}
	if m.InvestmentProduct != nil {
		dto.InvestmentProduct = new(e.InvestmentProductDto).FromModel(m.InvestmentProduct)
	}
	return dto
}

type UnitLinkFindAllRequest struct {
	entity.FindAllRequest
	ParticipantID string
}

func (r *UnitLinkFindAllRequest) GenerateFilter() {
	if r.ParticipantID != "" {
		r.AddFilter(pagination.FilterItem{
			Field: "participant_id",
			Op:    "eq",
			Val:   r.ParticipantID,
		})
	}
}

type UnitLinkSummaryPerTypeDto struct {
	Type                string  `json:"type"`
	SumIp               float64 `json:"sumIp"`
	TotalAmount         float64 `json:"totalAmount"`
	TotalAmountUnitLink float64 `json:"totalAmountUnitLink"`
}

type PortfolioWithNavDto struct {
	ID                  string                  `json:"id"`
	ParticipantID       string                  `json:"participantId"`
	Participant         *entity.ParticipantDto  `json:"participant"`
	CustomerID          string                  `json:"customerId"`
	Customer            *entity.CustomerDto     `json:"customer"`
	InvestmentProductID string                  `json:"investmentProductId"`
	InvestmentProduct   *e.InvestmentProductDto `json:"investmentProduct"`
	Ip                  float64                 `json:"ip"`
	LatestNav           float64                 `json:"latestNav"`
	TotalBalance        float64                 `json:"totalBalance"`
	TransactionDate     time.Time               `json:"transactionDate"`
	CreatedAt           time.Time               `json:"createdAt"`
}
