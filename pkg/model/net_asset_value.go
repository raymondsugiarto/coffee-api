package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type NetAssetValue struct {
	concern.CommonWithIDs
	OrganizationID      string
	Organization        *Organization
	InvestmentProductID string
	InvestmentProduct   *InvestmentProduct
	Amount              float64
	CreatedDate         time.Time
}
