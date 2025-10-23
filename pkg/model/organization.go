package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

// Accounts : table accounts
type Organization struct {
	concern.CommonWithIDs
	Code   string
	Name   string
	Origin string
}
