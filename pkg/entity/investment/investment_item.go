package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type InvestmentItemInputDto struct {
	InvestmentType      model.InvestmentType
	InvestmentProductID string
	Percent             float64
}

func (i *InvestmentItemInputDto) ToDto() *InvestmentItemDto {
	return &InvestmentItemDto{
		InvestmentType:      i.InvestmentType,
		InvestmentProductID: i.InvestmentProductID,
		Percent:             i.Percent,
	}
}

type InvestmentItemDto struct {
	ID                  string                 `json:"id,omitempty"`
	Code                string                 `json:"code,omitempty"`
	OrganizationID      string                 `json:"organization_id,omitempty"`
	ParticipantID       string                 `json:"participantId,omitempty"`
	Participant         *entity.ParticipantDto `json:"participant,omitempty"`
	CustomerID          string                 `json:"customerId,omitempty"`
	Customer            *entity.CustomerDto    `json:"customer,omitempty"`
	InvestmentID        *string                `json:"investmentId,omitempty"`
	Investment          *InvestmentDto         `json:"investment,omitempty"`
	InvestmentType      model.InvestmentType   `json:"investmentType,omitempty"`
	InvestmentProductID string                 `json:"investmentProductId,omitempty"`
	InvestmentProduct   *InvestmentProductDto  `json:"investmentProduct,omitempty"`
	Type                model.InvestmentType   `json:"type,omitempty"`
	Amount              float64                `json:"amount,omitempty"`      // Gross amount (net investment + fee) - total payment required
	FeeAmount           float64                `json:"feeAmount,omitempty"`   // Admin fee amount
	TotalAmount         float64                `json:"totalAmount,omitempty"` // Net investment amount (after fee deduction) - actual investment value
	EmployerAmount      float64                `json:"employerAmount,omitempty"`
	EmployeeAmount      float64                `json:"employeeAmount,omitempty"`
	VoluntaryAmount     float64                `json:"voluntaryAmount,omitempty"`
	EducationFundAmount float64                `json:"educationFundAmount,omitempty"`
	Percent             float64                `json:"percent,omitempty"`
	ExpiredAt           time.Time              `json:"expiredAt,omitempty"`
	InvestmentAt        time.Time              `json:"investmentAt,omitempty"`
	Status              model.InvestmentStatus `json:"status,omitempty"`
	CreatedAt           time.Time              `json:"createdAt,omitempty"`
}

func (d *InvestmentItemDto) ToModel() *model.InvestmentItem {
	m := &model.InvestmentItem{
		Code:                d.Code,
		OrganizationID:      d.OrganizationID,
		ParticipantID:       d.ParticipantID,
		CustomerID:          d.CustomerID,
		InvestmentType:      d.InvestmentType,
		InvestmentProductID: d.InvestmentProductID,
		Type:                d.Type,
		Amount:              d.Amount,
		FeeAmount:           d.FeeAmount,
		TotalAmount:         d.TotalAmount,
		EmployerAmount:      d.EmployerAmount,
		EmployeeAmount:      d.EmployeeAmount,
		VoluntaryAmount:     d.VoluntaryAmount,
		EducationFundAmount: d.EducationFundAmount,
		Percent:             d.Percent,
		ExpiredAt:           d.ExpiredAt,
		InvestmentAt:        d.InvestmentAt,
		Status:              d.Status,
	}

	if d.InvestmentID != nil {
		m.InvestmentID = *d.InvestmentID
	}

	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *InvestmentItemDto) FromModel(m *model.InvestmentItem) *InvestmentItemDto {
	d.ID = m.ID
	d.Code = m.Code
	d.OrganizationID = m.OrganizationID
	d.ParticipantID = m.ParticipantID
	d.CustomerID = m.CustomerID
	d.InvestmentType = m.InvestmentType
	d.InvestmentProductID = m.InvestmentProductID
	d.Amount = m.Amount
	d.FeeAmount = m.FeeAmount
	d.TotalAmount = m.TotalAmount
	d.EmployerAmount = m.EmployerAmount
	d.EmployeeAmount = m.EmployeeAmount
	d.VoluntaryAmount = m.VoluntaryAmount
	d.EducationFundAmount = m.EducationFundAmount
	d.Type = m.Type
	d.Percent = m.Percent
	d.ExpiredAt = m.ExpiredAt
	d.Status = m.Status
	d.CreatedAt = m.CreatedAt
	d.InvestmentAt = m.InvestmentAt

	if m.Investment != nil {
		d.Investment = (&InvestmentDto{}).FromModel(m.Investment)
	}

	if m.InvestmentProduct != nil {
		d.InvestmentProduct = (&InvestmentProductDto{}).FromModel(m.InvestmentProduct)
	}
	if m.Participant != nil {
		d.Participant = (&entity.ParticipantDto{}).FromModel(m.Participant)
	}
	if m.Customer != nil {
		d.Customer = (&entity.CustomerDto{}).FromModel(m.Customer)
	}
	return d
}

type InvestmentItemFindAllRequest struct {
	pagination.GetListRequest
	CustomerID       string
	CompanyID        string
	InvestmentID     string
	ShowAll          bool
	StartDate        time.Time
	EndDate          time.Time
	Status           model.InvestmentStatus
	InvestmentStatus model.InvestmentStatus
}

func (r *InvestmentItemFindAllRequest) GenerateFilter() {
	if !r.ShowAll {
		r.GetListRequest.AddFilter(pagination.FilterItem{
			Field: "investment_at",
			Op:    "gte",
			Val:   r.StartDate,
		})
		r.GetListRequest.AddFilter(pagination.FilterItem{
			Field: "investment_at",
			Op:    "lt",
			Val:   r.EndDate,
		})
	}
	if r.CustomerID != "" {
		r.GetListRequest.AddFilter(pagination.FilterItem{
			Field: "customer_id",
			Op:    "eq",
			Val:   r.CustomerID,
		})
	}
	if r.InvestmentID != "" {
		r.GetListRequest.AddFilter(pagination.FilterItem{
			Field: "investment_id",
			Op:    "eq",
			Val:   r.InvestmentID,
		})
	}

	if r.Status != "" {
		r.GetListRequest.AddFilter(pagination.FilterItem{
			Field: "status",
			Op:    "eq",
			Val:   r.Status,
		})
	}

}

type InvestmentStatementDto struct {
	*InvestmentItemDto
	Unit               float64   `json:"unit"`
	CurrentNavAmount   float64   `json:"currentNavAmount"`
	CurrentNavDate     time.Time `json:"currentNavDate"`
	CurrentValue       float64   `json:"currentValue"`
	GainLoss           float64   `json:"gainLoss"`
	GainLossPercentage float64   `json:"gainLossPercentage"`
}

func (i *InvestmentItemDto) ToInvestmentStatementDto() *InvestmentStatementDto {
	return &InvestmentStatementDto{
		InvestmentItemDto: i,
	}
}
