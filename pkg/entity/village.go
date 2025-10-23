package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type VillageDto struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Code       string       `json:"code"`
	PostalCode string       `json:"postal_code"`
	District   *DistrictDto `json:"district"`
}

func (d *VillageDto) ToModel() *model.Village {
	m := &model.Village{
		Name:       d.Name,
		Code:       d.Code,
		PostalCode: d.PostalCode,
	}
	if d.ID != "" {
		m.ID = d.ID
	}
	return m
}

func (d *VillageDto) FromModel(m *model.Village) *VillageDto {
	d.ID = m.ID
	d.Name = m.Name
	d.Code = m.Code
	d.PostalCode = m.PostalCode
	if m.District != nil {
		d.District = new(DistrictDto).FromModel(m.District)
	}
	return d
}

type VillageFindAllRequest struct {
	FindAllRequest
	Code         string
	DistrictID   string
	DistrictCode string
}

func (dto *VillageFindAllRequest) GenerateFilter() {
	if dto.Code != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "code",
				Op:    "eq",
				Val:   dto.Code,
			},
		)
	}
	if dto.DistrictID != "" {
		dto.AddFilter(
			pagination.FilterItem{
				Field: "district_id",
				Op:    "eq",
				Val:   dto.DistrictID,
			},
		)
	}
}
