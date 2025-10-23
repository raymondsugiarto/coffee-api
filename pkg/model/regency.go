package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Regency struct {
	concern.CommonWithIDs
	Name       string
	Code       string
	ProvinceID string
	Province   *Province
}
