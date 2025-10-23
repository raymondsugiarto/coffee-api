package model

import "github.com/raymondsugiarto/coffee-api/pkg/model/concern"

type Country struct {
	concern.CommonWithIDs
	Name string
	CCA2 string
	CCA3 string
}
