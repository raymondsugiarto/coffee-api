package participantsummary

import (
	"context"
	"errors"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionPerCompanyTypeAndPilar(ctx context.Context, filter *entity.ReportSummaryAumFilter) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetTransactionPerCompanyTypeAndPilar(ctx context.Context, filter *entity.ReportSummaryAumFilter) (*pagination.ResultPagination, error) {
	var m []entity.ReportSummaryAum = make([]entity.ReportSummaryAum, 0)

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	startDate := fmt.Sprintf("%d-%02d-01", filter.Year, filter.Month)
	endMonth := filter.Month + 1
	endYear := filter.Year
	if endMonth > 12 {
		endMonth = 1
		endYear++
	}
	endDate := fmt.Sprintf("%d-%02d-01", endYear, endMonth)
	today := time.Now().Format("2006-01-02")

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		db := r.db.WithContext(ctx)

		subQuery := db.Table("unit_link").
			Select("unit_link.participant_id, unit_link.investment_product_id, SUM(unit_link.ip * COALESCE(net_asset_value.amount, 0)) AS aum").
			Joins("LEFT JOIN net_asset_value ON net_asset_value.investment_product_id = unit_link.investment_product_id AND net_asset_value.created_date = ?", today).
			Where("unit_link.deleted_at IS NULL").
			Where("unit_link.transaction_date >= ?", startDate).
			Where("unit_link.transaction_date < ?", endDate).
			Group("unit_link.participant_id, unit_link.investment_product_id")

		q := db.Table("(?) as a", subQuery).
			Select([]string{
				"company.company_type",
				"company.pilar_type",
				"COALESCE(SUM(a.aum), 0) AS aum",
			}).
			Joins("JOIN participant ON participant.id = a.participant_id").
			Joins("JOIN customer ON customer.id = participant.customer_id").
			Joins("JOIN company ON company.id = customer.company_id").
			Group("company.company_type, company.pilar_type")

		return q
	}, &pagination.TableRequest{
		Request:       filter,
		QueryField:    []string{"company.company_type", "company.pilar_type"},
		Data:          &m,
		AllowedFields: []string{"company_type", "pilar_type", "aum"},
		MapFields: map[string]string{
			"company_type": "company.company_type",
			"pilar_type":   "company.pilar_type",
			"aum":          "aum",
		},
	})

	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]entity.ReportSummaryAum)
	var data []*entity.ReportSummaryAum = make([]*entity.ReportSummaryAum, 0)

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
