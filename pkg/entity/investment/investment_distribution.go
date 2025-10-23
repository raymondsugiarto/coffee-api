package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type InvestmentDistributionInputBatchDto struct {
	Data []*InvestmentDistributionInputDto `json:"data"`
}

type InvestmentDistributionInputDto struct {
	ID                    *string                `json:"id,omitempty"`
	Type                  string                 `json:"type,omitempty"` // fill this with company or participant
	CompanyID             string                 `json:"companyId,omitempty"`
	Company               *entity.CompanyDto     `json:"company,omitempty"`
	ParticipantID         string                 `json:"participantId,omitempty"`
	Participant           *entity.ParticipantDto `json:"participant,omitempty"`
	CustomerID            string                 `json:"customerId,omitempty"`
	Customer              *entity.CustomerDto    `json:"customer,omitempty"`
	InvestmentProductID   string                 `json:"investmentProductId,omitempty"`
	InvestmentProduct     *InvestmentProductDto  `json:"investmentProduct,omitempty"`
	Percent               float64                `json:"percent"`
	BaseContribution      float64                `json:"baseContribution,omitempty"`
	VoluntaryContribution float64                `json:"voluntaryContribution,omitempty"`
}

func (i *InvestmentDistributionInputDto) ToDto() *InvestmentDistributionDto {
	return &InvestmentDistributionDto{
		ID:                    *i.ID,
		Type:                  i.Type,
		CompanyID:             i.CompanyID,
		Company:               i.Company,
		ParticipantID:         i.ParticipantID,
		Participant:           i.Participant,
		CustomerID:            i.CustomerID,
		Customer:              i.Customer,
		InvestmentProductID:   i.InvestmentProductID,
		InvestmentProduct:     i.InvestmentProduct,
		Percent:               i.Percent,
		BaseContribution:      i.BaseContribution,
		VoluntaryContribution: i.VoluntaryContribution,
	}
}

type InvestmentDistributionDto struct {
	ID                    string                 `json:"id,omitempty"`
	OrganizationID        string                 `json:"organizationId,omitempty"`
	Type                  string                 `json:"type,omitempty"`
	CompanyID             string                 `json:"companyId,omitempty"`
	Company               *entity.CompanyDto     `json:"company,omitempty"`
	ParticipantID         string                 `json:"participantId,omitempty"`
	Participant           *entity.ParticipantDto `json:"participant,omitempty"`
	CustomerID            string                 `json:"customerId,omitempty"`
	Customer              *entity.CustomerDto    `json:"customer,omitempty"`
	InvestmentProductID   string                 `json:"investmentProductId,omitempty"`
	InvestmentProduct     *InvestmentProductDto  `json:"investmentProduct,omitempty"`
	Percent               float64                `json:"percent"`
	BaseContribution      float64                `json:"baseContribution,omitempty"`
	VoluntaryContribution float64                `json:"voluntaryContribution,omitempty"`
}

func (d *InvestmentDistributionDto) ToUserLogDto() *entity.UserLogDto {
	return &entity.UserLogDto{
		OrganizationID: d.OrganizationID,
		CompanyID:      d.CompanyID,
		RefModule:      "investment_distribution",
		RefTable:       "investment_distribution",
		RefID:          d.ID,
	}
}

func (d *InvestmentDistributionDto) ToModel() *model.InvestmentDistribution {
	m := &model.InvestmentDistribution{
		OrganizationID:        d.OrganizationID,
		Type:                  d.Type,
		CompanyID:             d.CompanyID,
		ParticipantID:         d.ParticipantID,
		CustomerID:            d.CustomerID,
		InvestmentProductID:   d.InvestmentProductID,
		Percent:               d.Percent,
		BaseContribution:      d.BaseContribution,
		VoluntaryContribution: d.VoluntaryContribution,
	}

	if d.Company != nil {
		m.Company = d.Company.ToModel()
	}
	if d.Participant != nil {
		m.Participant = d.Participant.ToModel()
	}
	if d.Customer != nil {
		m.Customer = d.Customer.ToModel()
	}
	if d.InvestmentProduct != nil {
		m.InvestmentProduct = d.InvestmentProduct.ToModel()
	}
	return m
}

func (d *InvestmentDistributionDto) FromModel(m *model.InvestmentDistribution) *InvestmentDistributionDto {
	d.ID = m.ID
	d.Type = m.Type
	d.CompanyID = m.CompanyID
	d.ParticipantID = m.ParticipantID
	d.CustomerID = m.CustomerID
	d.InvestmentProductID = m.InvestmentProductID
	d.Percent = m.Percent
	d.BaseContribution = m.BaseContribution
	d.VoluntaryContribution = m.VoluntaryContribution
	if m.Company != nil {
		d.Company = (&entity.CompanyDto{}).FromModel(m.Company)
	}
	if m.Participant != nil {
		d.Participant = (&entity.ParticipantDto{}).FromModel(m.Participant)
	}
	if m.Customer != nil {
		d.Customer = (&entity.CustomerDto{}).FromModel(m.Customer)
	}
	if m.InvestmentProduct != nil {
		d.InvestmentProduct = (&InvestmentProductDto{}).FromModel(m.InvestmentProduct)
	}
	return d
}

type InvestmentDistributionSummaryCompanyDto struct {
	InvestmentProductID   string  `json:"investmentProductId"`
	InvestmentProductName string  `json:"investmentProductName"`
	InvestAmount          float64 `json:"investAmount"`
	TotalAmount           float64 `json:"totalAmount"`
	Ip                    float64 `json:"ip"`
	Nab                   float64 `json:"nab"`
}
