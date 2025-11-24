package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
)

type CompanyDto struct {
	ID             string
	OrganizationID string
	Organization   OrganizationDto
	PhoneNumber    string
	Name           string
}

func NewCompanyDtoFromModel(m *model.Company) *CompanyDto {
	if m == nil {
		return nil
	}
	d := &CompanyDto{
		ID:   m.ID,
		Name: m.Name,
	}
	return d
}
