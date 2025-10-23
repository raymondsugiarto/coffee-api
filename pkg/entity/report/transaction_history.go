package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type TransactionHistoryFilter struct {
	entity.FindAllRequest
	CompanyID string                 `query:"companyId"`
	StartDate time.Time              `query:"startDate"`
	EndDate   time.Time              `query:"endDate"`
	Status    model.InvestmentStatus `query:"status"`
}

func (f *TransactionHistoryFilter) GenerateFilter() {
	if f.CompanyID != "" {
		f.AddFilter(pagination.FilterItem{
			Field: "company_id",
			Op:    "eq",
			Val:   f.CompanyID,
		})
	}

	if f.Status != "" {
		f.AddFilter(pagination.FilterItem{
			Field: "status",
			Op:    "eq",
			Val:   f.Status,
		})
	}

	if !f.StartDate.IsZero() {
		f.AddFilter(pagination.FilterItem{
			Field: "investment_at",
			Op:    "gte",
			Val:   f.StartDate,
		})
	}

	if !f.EndDate.IsZero() {
		f.AddFilter(pagination.FilterItem{
			Field: "investment_at",
			Op:    "lte",
			Val:   f.EndDate,
		})
	}
}

type TransactionHistoryReport struct {
	InvestmentItemCode    string    `json:"investmentItemCode"`
	ParticipantCode       string    `json:"participantCode"`
	ParticipantName       string    `json:"participantName"`
	CompanyCode           string    `json:"companyCode"`
	CompanyName           string    `json:"companyName"`
	InvestmentProductCode string    `json:"investmentProductCode"`
	InvestmentProductName string    `json:"investmentProductName"`
	Amount                float64   `json:"amount"`
	FeeAmount             float64   `json:"feeAmount"`
	TotalAmount           float64   `json:"totalAmount"`
	NavAmount             *float64  `json:"navAmount"`
	UnitAmount            float64   `json:"unitAmount"`
	InvestmentAt          time.Time `json:"investmentAt"`
	Status                string    `json:"status"`
}
