package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type PensionBenefitRecipientInputDto struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Relationship         string `json:"relationship,omitempty"`
	DateOfBirth          string `json:"dateOfBirth,omitempty"`
	CountryOfBirth       string `json:"countryOfBirth,omitempty"`
	IdentificationNumber string `json:"identificationNumber,omitempty"`
	PhoneNumber          string `json:"phoneNumber,omitempty"`
	CustomerID           string `json:"customerId,omitempty"`
}

func (i *PensionBenefitRecipientInputDto) ToDto() *PensionBenefitRecipientDto {
	dto := &PensionBenefitRecipientDto{
		Name:                 i.Name,
		Relationship:         i.Relationship,
		DateOfBirth:          i.DateOfBirth,
		CountryOfBirth:       i.CountryOfBirth,
		IdentificationNumber: i.IdentificationNumber,
		PhoneNumber:          i.PhoneNumber,
		CustomerID:           i.CustomerID,
	}
	if i.ID != "" {
		dto.ID = i.ID
	}
	return dto
}

type PensionBenefitRecipientDto struct {
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	Relationship         string      `json:"relationship"`
	DateOfBirth          string      `json:"dateOfBirth"`
	CountryOfBirth       string      `json:"countryOfBirth"`
	CountryBirth         *CountryDto `json:"countryBirth"`
	IdentificationNumber string      `json:"identificationNumber"`
	PhoneNumber          string      `json:"phoneNumber"`
	CustomerID           string      `json:"customerId"`
}

func (d *PensionBenefitRecipientDto) ToModel() *model.PensionBenefitRecipient {
	m := &model.PensionBenefitRecipient{
		Name:                 d.Name,
		Relationship:         d.Relationship,
		DateOfBirth:          d.DateOfBirth,
		CountryOfBirth:       d.CountryOfBirth,
		IdentificationNumber: d.IdentificationNumber,
		PhoneNumber:          d.PhoneNumber,
		CustomerID:           d.CustomerID,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *PensionBenefitRecipientDto) FromModel(m *model.PensionBenefitRecipient) *PensionBenefitRecipientDto {
	d.ID = m.ID
	d.Name = m.Name
	d.Relationship = m.Relationship
	d.DateOfBirth = m.DateOfBirth
	d.CountryOfBirth = m.CountryOfBirth
	d.IdentificationNumber = m.IdentificationNumber
	d.PhoneNumber = m.PhoneNumber
	d.CustomerID = m.CustomerID
	if m.CountryBirth != nil {
		d.CountryBirth = (&CountryDto{}).FromModel(m.CountryBirth)
	}
	return d
}

type RecipientFindAllRequest struct {
	FindAllRequest
	CustomerID string
}

func (dto *RecipientFindAllRequest) GenerateFilter() {
	if dto.CustomerID != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "customer_id",
				Op:    "eq",
				Val:   dto.CustomerID,
			},
		)
	}
}
