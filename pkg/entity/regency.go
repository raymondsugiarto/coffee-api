package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type RegencyDto struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Code     string       `json:"code"`
	Province *ProvinceDto `json:"province"`
}

func (d *RegencyDto) ToModel() *model.Regency {
	m := &model.Regency{
		Name: d.Name,
		Code: d.Code,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *RegencyDto) FromModel(m *model.Regency) *RegencyDto {
	d.ID = m.ID
	d.Name = m.Name
	d.Code = m.Code
	if m.Province != nil {
		d.Province = new(ProvinceDto).FromModel(m.Province)
	}
	return d
}

type RegencyFindAllRequest struct {
	FindAllRequest
	Code         string
	ProvinceID   string
	ProvinceCode string
}

func (dto *RegencyFindAllRequest) GenerateFilter() {
	if dto.Code != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "code",
				Op:    "eq",
				Val:   dto.Code,
			},
		)
	}
	if dto.ProvinceID != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "province_id",
				Op:    "eq",
				Val:   dto.ProvinceID,
			},
		)
	}
}
