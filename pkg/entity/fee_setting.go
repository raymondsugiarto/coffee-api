package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type FeeSettingDto struct {
	ID             string  `json:"id"`
	AdminFee       float64 `json:"adminFee"`
	OperationalFee float64 `json:"operationalFee"`
}

func (dto *FeeSettingDto) ToModel() *model.FeeSetting {
	return &model.FeeSetting{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
		AdminFee:       dto.AdminFee,
		OperationalFee: dto.OperationalFee,
	}
}

func (dto *FeeSettingDto) FromModel(m *model.FeeSetting) *FeeSettingDto {
	return &FeeSettingDto{
		ID:             m.ID,
		AdminFee:       m.AdminFee,
		OperationalFee: m.OperationalFee,
	}
}
