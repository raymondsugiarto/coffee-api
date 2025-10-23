package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	ei "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type PortfolioDto struct {
	ID                  string                   `json:"id"`
	OrganizationID      string                   `json:"organizationId"`
	Organization        *entity.OrganizationData `json:"-"`
	CompanyID           string                   `json:"companyId"`
	Company             *entity.CompanyDto       `json:"company"`
	CustomerID          string                   `json:"customerId"`
	Customer            *entity.CustomerDto      `json:"customer"`
	ParticipantID       string                   `json:"participantId"`
	Participant         *entity.ParticipantDto   `json:"participant"`
	InvestmentProductID string                   `json:"investmentProductId"`
	InvestmentProduct   *ei.InvestmentProductDto `json:"investmentProduct"`
	Ip                  float64                  `json:"ip"`
}

func (p *PortfolioDto) ToModel() *model.Portfolio {
	r := &model.Portfolio{
		OrganizationID:      p.OrganizationID,
		CompanyID:           p.CompanyID,
		CustomerID:          p.CustomerID,
		ParticipantID:       p.ParticipantID,
		InvestmentProductID: p.InvestmentProductID,
		Ip:                  p.Ip,
	}
	if p.ID != "" {
		r.ID = p.ID
	}
	return r
}

func (p *PortfolioDto) FromModel(m *model.Portfolio) *PortfolioDto {
	if m == nil {
		return nil
	}
	p.ID = m.ID
	p.OrganizationID = m.OrganizationID
	p.CompanyID = m.CompanyID
	p.CustomerID = m.CustomerID
	p.ParticipantID = m.ParticipantID
	p.InvestmentProductID = m.InvestmentProductID
	p.Ip = m.Ip
	if m.Company != nil {
		p.Company = new(entity.CompanyDto).FromModel(m.Company)
	}
	if m.Customer != nil {
		p.Customer = new(entity.CustomerDto).FromModel(m.Customer)
	}
	if m.Participant != nil {
		p.Participant = new(entity.ParticipantDto).FromModel(m.Participant)
	}
	if m.InvestmentProduct != nil {
		p.InvestmentProduct = new(ei.InvestmentProductDto).FromModel(m.InvestmentProduct)
	}
	return p
}

type PortfolioFindAllRequest struct {
	entity.FindAllRequest
	CustomerID    string `json:"customerId"`
	ParticipantID string `json:"participantId"`
}
