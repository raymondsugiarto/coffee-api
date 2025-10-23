package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type BenefitType struct {
	concern.CommonWithIDs

	Name                    string
	Description             string
	MinimumTimePeriodMonths int
	MinimumContribution     float64
	Status                  BenefitTypeStatus
}

type BenefitTypeStatus string

const (
	ACTIVE   BenefitTypeStatus = "ACTIVE"
	INACTIVE BenefitTypeStatus = "INACTIVE"
)
