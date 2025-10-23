package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Village struct {
	concern.CommonWithIDs
	Name       string
	Code       string
	PostalCode string
	DistrictID string
	District   *District
}
