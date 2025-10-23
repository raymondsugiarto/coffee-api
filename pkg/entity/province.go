package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ProvinceDto struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

func (d *ProvinceDto) ToModel() *model.Province {
	m := &model.Province{
		Name: d.Name,
		Code: d.Code,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *ProvinceDto) FromModel(m *model.Province) *ProvinceDto {
	return &ProvinceDto{
		ID:   m.ID,
		Name: m.Name,
		Code: m.Code,
	}
}

type ProvinceFindAllRequest struct {
	FindAllRequest
	Code string
}

func (dto *ProvinceFindAllRequest) GenerateFilter() {
	if dto.Code != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "code",
				Op:    "eq",
				Val:   dto.Code,
			},
		)
	}
}
