package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Province struct {
	concern.CommonWithIDs
	Name string
	Code string
}
