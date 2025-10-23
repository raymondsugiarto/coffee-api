package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type DistrictDto struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Code    string      `json:"code"`
	Regency *RegencyDto `json:"regency"`
}

func (d *DistrictDto) ToModel() *model.District {
	m := &model.District{
		Name: d.Name,
		Code: d.Code,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *DistrictDto) FromModel(m *model.District) *DistrictDto {
	d.ID = m.ID
	d.Name = m.Name
	d.Code = m.Code
	if m.Regency != nil {
		d.Regency = new(RegencyDto).FromModel(m.Regency)
	}
	return d
}

type DistrictFindAllRequest struct {
	FindAllRequest
	Code        string
	RegencyID   string
	RegencyCode string
}

func (dto *DistrictFindAllRequest) GenerateFilter() {
	if dto.Code != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "code",
				Op:    "eq",
				Val:   dto.Code,
			},
		)
	}
	if dto.RegencyID != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "regency_id",
				Op:    "eq",
				Val:   dto.RegencyID,
			},
		)
	}
}
