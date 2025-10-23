package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type ReportSummaryAumFilter struct {
	entity.FindAllRequest
	Month int64
	Year  int64
}

type ReportSummaryAum struct {
	CompanyType string
	PilarType   string
	Aum         float64
}

type ReportSummaryAumAggregated struct {
	GroupProduksi string  `json:"groupProduksi"` // PILAR / NON PILAR
	AumPPIP       float64 `json:"aumPPIP"`
	AumDKP        float64 `json:"aumDKP"`
	TotalAum      float64 `json:"totalAum"`
}
