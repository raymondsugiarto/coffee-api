package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ReportTransactionChannel struct {
	ParticipantID         string     `json:"participantID"`
	ParticipantName       string     `json:"participantName"`
	InvestmentProductID   string     `json:"investmentProductID"`
	InvestmentProductName string     `json:"investmentProductName"`
	CompanyName           string     `json:"companyName"`
	Balance               float64    `json:"balance"`
	Nab                   float64    `json:"nab"`
	Fee                   float64    `json:"fee"`
	InvestmentAt          *time.Time `json:"investmentAt"`
}

type ReportTransactionChannelFilter struct {
	entity.FindAllRequest

	EndDate *time.Time
}

func (f *ReportTransactionChannelFilter) GenerateFilter() {
	if !f.EndDate.IsZero() {
		f.AddFilter(pagination.FilterItem{
			Field: "investment_at",
			Op:    "lte",
			Val:   f.EndDate.Add(24*time.Hour - time.Second),
		})
	}
}
