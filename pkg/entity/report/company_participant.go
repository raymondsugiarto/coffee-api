package entity

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type CompanyParticipantFilter struct {
	entity.FindAllRequest
	CompanyID string `query:"companyId"`
}

func (f *CompanyParticipantFilter) GenerateFilter() {
	if f.CompanyID != "" {
		f.AddFilter(pagination.FilterItem{
			Field: "company.id",
			Op:    "eq",
			Val:   f.CompanyID,
		})
	}
}

type CompanyParticipantReport struct {
	ParticipantID     string  `json:"participantId"`
	CustomerID        string  `json:"customerId"`
	CustomerName      string  `json:"customerName"`
	ParticipantCode   string  `json:"participantCode"`
	TotalContribution float64 `json:"totalContribution"`
	TotalUnit         float64 `json:"totalUnit"`
	LastBalance       float64 `json:"lastBalance"`
}
