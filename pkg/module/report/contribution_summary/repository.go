package contributionsummary

import (
	"context"
	"errors"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	GetReportContribution(ctx context.Context, filter *entity.ReportContributionSummaryFilter) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetReportContribution(ctx context.Context, filter *entity.ReportContributionSummaryFilter) (*pagination.ResultPagination, error) {
	var m []entity.ReportContributionSummary = make([]entity.ReportContributionSummary, 0)

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		db := r.db.WithContext(ctx)

		subQuery := db.Model(&model.Customer{}).
			Select([]string{
				"customer.company_id",
				"COALESCE(SUM(customer.customer_amount), 0) AS customer_amount",
				"COALESCE(SUM(customer.voluntary_amount), 0) AS voluntary_amount",
				"COALESCE(SUM(customer.employer_amount), 0) AS employer_amount",
				"COALESCE(SUM(customer.education_fund_amount), 0) AS education_fund_amount",
			}).
			Where("customer.sim_status = ?", model.SIMStatusActive)

		if filter.EndDate != nil {
			subQuery = subQuery.Where("customer.effective_date <= ?", filter.EndDate)
		}

		subQuery = subQuery.Group("customer.company_id")

		q := db.Table("(?) AS c", subQuery).
			Select([]string{
				`COALESCE(company.first_name, 'PESERTA MANDIRI') AS name`,
				"c.customer_amount",
				"c.voluntary_amount",
				"c.employer_amount",
				"c.education_fund_amount",
				"(c.customer_amount + c.voluntary_amount + c.employer_amount + c.education_fund_amount) AS total",
				`CASE 
					WHEN company.company_type = 'PPIP' OR NULLIF(TRIM(company.id), '') IS NULL 
					THEN '1001' ELSE '1002' 
				END AS type_code`,
			}).
			Joins("LEFT JOIN company ON company.id = c.company_id").
			Where(`
				NULLIF(TRIM(c.company_id), '') IS NULL
				OR (company.status = 'APPROVED' AND company.deleted_at IS NULL)
			`).
			Order("CASE WHEN company.first_name IS NULL THEN 0 ELSE 1 END, company.first_name")

		if filter.EndDate != nil {
			q = q.Where("NULLIF(TRIM(c.company_id), '') IS NULL OR company.created_at <= ?", filter.EndDate)
		}

		return q
	}, &pagination.TableRequest{
		Request:       filter,
		QueryField:    []string{"company.first_name"},
		Data:          &m,
		AllowedFields: []string{"name", "customer_amount", "voluntary_amount", "employer_amount", "total", "type_code", "company.created_at", "education_fund_amount"},
		MapFields: map[string]string{
			"name":                  "company.first_name",
			"customer_amount":       "c.customer_amount",
			"voluntary_amount":      "c.voluntary_amount",
			"employer_amount":       "c.employer_amount",
			"total":                 "total",
			"type_code":             "type_code",
			"company.created_at":    "company.created_at",
			"education_fund_amount": "c.education_fund_amount",
		},
	})

	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]entity.ReportContributionSummary)
	var data []*entity.ReportContributionSummary = make([]*entity.ReportContributionSummary, 0)

	// Copy data
	for _, item := range *results {
		itemCopy := item
		data = append(data, &itemCopy)
	}

	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
