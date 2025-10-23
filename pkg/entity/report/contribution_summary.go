package entity

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type ReportContributionSummary struct {
	Name                string  `json:"name"`
	CustomerAmount      float64 `json:"customerAmount"`
	VoluntaryAmount     float64 `json:"voluntaryAmount"`
	EmployerAmount      float64 `json:"employerAmount"`
	Total               float64 `json:"total"`
	TypeCode            string  `json:"typeCode"`
	EducationFundAmount float64 `json:"educationFundAmount"`
}

type ReportContributionSummaryFilter struct {
	entity.FindAllRequest
	EndDate *time.Time
}

func (f *ReportContributionSummaryFilter) GenerateFilter() {
	if f.EndDate != nil && !f.EndDate.IsZero() {
		f.AddFilter(pagination.FilterItem{
			Field: "company.created_at",
			Op:    "lte",
			Val:   f.EndDate,
		})
	}
}
