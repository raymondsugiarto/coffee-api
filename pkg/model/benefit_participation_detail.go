package model

import (
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type BenefitParticipationDetail struct {
	concern.CommonWithIDs

	BenefitParticipationID  string
	BenefitTypeID           string
	BenefitType             BenefitType
	TimePeriodMonths        int
	PlannedWithdrawalMonths int
	MonthlyContribution     float64
	Status                  BenefitParticipationStatus
}
