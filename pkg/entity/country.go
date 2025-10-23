package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type CountryDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CCA2 string `json:"cca2"`
	CCA3 string `json:"cca3"`
}

func (d *CountryDto) ToModel() *model.Country {
	m := &model.Country{
		Name: d.Name,
		CCA2: d.CCA2,
		CCA3: d.CCA3,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *CountryDto) FromModel(m *model.Country) *CountryDto {
	return &CountryDto{
		ID:   m.ID,
		Name: m.Name,
		CCA2: m.CCA2,
		CCA3: m.CCA3,
	}
}

type CountryFindAllRequest struct {
	FindAllRequest
	CCA2 string
	CCA3 string
}

func (dto *CountryFindAllRequest) GenerateFilter() {
	if dto.CCA2 != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "cca2",
				Op:    "eq",
				Val:   dto.CCA2,
			},
		)
	}
	if dto.CCA3 != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "cca3",
				Op:    "eq",
				Val:   dto.CCA3,
			},
		)
	}
}
