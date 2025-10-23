package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type FeeSetting struct {
	concern.CommonWithIDs

	AdminFee       float64
	OperationalFee float64
}
