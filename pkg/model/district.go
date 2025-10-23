package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type District struct {
	concern.CommonWithIDs
	Name      string
	Code      string
	RegencyID string
	Regency   *Regency
}
