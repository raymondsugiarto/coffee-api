package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type InvestmentInputDto struct {
	CompanyID        string // if not empty, will be used on company site
	CustomerID       string
	ParticipantID    string // if not empty, will be used on mobile site
	InvestmentAt     time.Time
	Amount           float64 `validate:"required,gt=0"`
	IsNewParticipant bool

	// TODO: bukti bayar
	InvestmentItems []*InvestmentItemInputDto
}

func (i *InvestmentInputDto) ToDto() *InvestmentDto {
	items := make([]*InvestmentItemDto, len(i.InvestmentItems))
	for j, item := range i.InvestmentItems {
		itemDto := item.ToDto()
		if !i.IsNewParticipant {
			itemDto.ParticipantID = i.ParticipantID
			itemDto.CustomerID = i.CustomerID
		}
		items[j] = itemDto
	}

	return &InvestmentDto{
		CompanyID:        i.CompanyID,
		CustomerID:       i.CustomerID,
		ParticipantID:    i.ParticipantID,
		InvestmentAt:     i.InvestmentAt,
		Amount:           i.Amount,
		IsNewParticipant: i.IsNewParticipant,
		InvestmentItems:  items,
	}
}

type InvestmentDto struct {
	ID                 string                  `json:"id,omitempty"`
	Code               string                  `json:"code,omitempty"`
	OrganizationID     string                  `json:"organization_id,omitempty"`
	CompanyID          string                  `json:"companyId,omitempty"`
	Company            *entity.CompanyDto      `json:"company,omitempty"`
	ParticipantID      string                  `json:"participantId,omitempty"`
	Participant        *entity.ParticipantDto  `json:"participant,omitempty"`
	CustomerID         string                  `json:"customerId,omitempty"`
	Customer           *entity.CustomerDto     `json:"customer,omitempty"`
	Amount             float64                 `json:"amount,omitempty"`
	Type               model.InvestmentType    `json:"type,omitempty"`
	ExpiredAt          time.Time               `json:"expiredAt,omitempty"`
	InvestmentAt       time.Time               `json:"investmentAt,omitempty"`
	Status             model.InvestmentStatus  `json:"status,omitempty"`
	Source             model.InvestmentSource  `json:"source,omitempty"`
	CreatedAt          time.Time               `json:"createdAt,omitempty"`
	IsNewParticipant   bool                    `json:"-"`
	InvestmentPayments []*InvestmentPaymentDto `json:"investmentPayments,omitempty"`
	InvestmentItems    []*InvestmentItemDto    `json:"investmentItems,omitempty"`
}

func (d *InvestmentDto) ToInvestmentDistributions() []*InvestmentDistributionDto {
	items := make([]*InvestmentDistributionDto, len(d.InvestmentItems))
	for j, item := range d.InvestmentItems {
		itemDto := new(InvestmentDistributionDto)
		itemDto.ParticipantID = d.ParticipantID
		itemDto.CustomerID = d.CustomerID
		itemDto.InvestmentProductID = item.InvestmentProductID
		itemDto.Percent = item.Percent
		itemDto.BaseContribution = d.Amount * item.Percent / 100
		items[j] = itemDto
	}
	return items
}

func (d *InvestmentDto) GetTotalPaymentAmount() float64 {
	// Calculate total payment amount from all investment items (gross amount including fees)
	totalPaymentAmount := 0.0
	for _, item := range d.InvestmentItems {
		totalPaymentAmount += item.Amount
	}
	return totalPaymentAmount
}

func (d *InvestmentDto) ToCompanyInvestmentPayment() *InvestmentPaymentDto {
	return &InvestmentPaymentDto{
		InvestmentID: d.ID,
		PaymentAt:    d.InvestmentAt,
		Amount:       d.GetTotalPaymentAmount(),
		Status:       model.InvestmentPaymentStatusConfirmed,
	}

}

func (d *InvestmentDto) ToModel() *model.Investment {
	m := &model.Investment{
		OrganizationID: d.OrganizationID,
		CompanyID:      d.CompanyID,
		ParticipantID:  d.ParticipantID,
		CustomerID:     d.CustomerID,
		Type:           d.Type,
		Amount:         d.Amount,
		ExpiredAt:      d.ExpiredAt,
		InvestmentAt:   d.InvestmentAt,
		Status:         d.Status,
		Code:           d.Code,
		Source:         d.Source,
	}

	investmentItems := make([]*model.InvestmentItem, 0)
	for _, item := range d.InvestmentItems {
		investmentItems = append(investmentItems, item.ToModel())
	}

	m.InvestmentItems = investmentItems

	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *InvestmentDto) FromModel(m *model.Investment) *InvestmentDto {
	d.ID = m.ID
	d.OrganizationID = m.OrganizationID
	d.CompanyID = m.CompanyID
	d.CustomerID = m.CustomerID
	d.ParticipantID = m.ParticipantID
	d.Amount = m.Amount
	d.ExpiredAt = m.ExpiredAt
	d.Type = m.Type
	d.Status = m.Status
	d.Source = m.Source
	d.CreatedAt = m.CreatedAt
	d.InvestmentAt = m.InvestmentAt
	d.Code = m.Code

	if m.Company != nil {
		d.Company = (&entity.CompanyDto{}).FromModel(m.Company)
	}
	if m.Customer != nil {
		d.Customer = (&entity.CustomerDto{}).FromModel(m.Customer)
	}
	if m.Participant != nil {
		d.Participant = (&entity.ParticipantDto{}).FromModel(m.Participant)
	}
	if m.InvestmentItems != nil {
		d.InvestmentItems = make([]*InvestmentItemDto, len(m.InvestmentItems))
		for i, v := range m.InvestmentItems {
			d.InvestmentItems[i] = (&InvestmentItemDto{}).FromModel(v)
		}
	}
	if m.InvestmentPayments != nil {
		d.InvestmentPayments = make([]*InvestmentPaymentDto, len(m.InvestmentPayments))
		for i, v := range m.InvestmentPayments {
			d.InvestmentPayments[i] = (&InvestmentPaymentDto{}).FromModel(v)
		}
	}
	return d

}

type InvestmentFindAllRequest struct {
	entity.FindAllRequest
	CustomerID      string
	ShowAll         bool
	StartDate       time.Time
	EndDate         time.Time
	CompanyID       *string
	IncludePayments bool
}

func (r *InvestmentFindAllRequest) GenerateFilter() {
	if r.CustomerID != "" {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "customer_id",
			Op:    "eq",
			Val:   r.CustomerID,
		})
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "type",
			Op:    "eq",
			Val:   "PPIP",
		})
	}
	if !r.ShowAll {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "created_at",
			Op:    "gte",
			Val:   r.StartDate,
		})
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "created_at",
			Op:    "lt",
			Val:   r.EndDate,
		})
	}
	if r.CompanyID != nil {
		r.FindAllRequest.AddFilter(pagination.FilterItem{
			Field: "company_id",
			Op:    "eq",
			Val:   r.CompanyID,
		})
	}
}

func (i *InvestmentDto) GetInfo() entity.RejectEmail {
	return entity.RejectEmail{
		Email:       i.Company.Email,
		Name:        i.Code,
		Description: "Pengajuan Investasi Dengan Kode " + i.Code,
	}
}
