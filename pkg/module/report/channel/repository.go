package channel

import (
	"context"
	"errors"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetTransactionReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) (*pagination.ResultPagination, error) {
	var m []entity.ReportTransactionChannel = make([]entity.ReportTransactionChannel, 0)

	if filter == nil {
		return nil, errors.New("filter is nil")
	}

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		subquery := r.db.WithContext(ctx).Table("investment_item").
			Select("investment_item.participant_id, investment_item.investment_product_id, SUM(investment_item.amount / net_asset_value.amount) AS nab, SUM(investment_item.fee_amount) AS fee, MAX(investment_item.investment_at) as investment_at").
			Joins("JOIN investment ON investment_item.investment_id = investment.id").
			Joins("LEFT JOIN net_asset_value ON investment_item.investment_product_id = net_asset_value.investment_product_id AND net_asset_value.created_date = date(investment_item.investment_at)").
			Where("investment.status = ? AND investment_item.deleted_at is null AND investment.deleted_at is null", "SUCCESS").
			Where("investment_item.investment_at < ?", filter.EndDate.Add(24*time.Hour)).
			Group("investment_item.participant_id, investment_item.investment_product_id")

		q := r.db.WithContext(ctx).Table("(?) as a", subquery).
			Select(`a.investment_product_id,participant.code as participant_id,customer.first_name as participant_name,company.first_name as company_name,investment_product.name as investment_product_name,a.nab,a.fee
			`).
			Joins("LEFT JOIN participant ON participant.id = a.participant_id").
			Joins("LEFT JOIN customer ON customer.id = participant.customer_id").
			Joins("LEFT JOIN company ON company.id = customer.company_id").
			Joins("LEFT JOIN investment_product ON investment_product.id = a.investment_product_id")

		return q
	}, &pagination.TableRequest{
		Request:       filter,
		QueryField:    []string{"participant.code", "customer.first_name", "company.first_name", "investment_at"},
		Data:          &m,
		AllowedFields: []string{"participant_id", "participant_name", "company_name", "investment_product_name", "nab", "fee", "investment_at"},
		MapFields: map[string]string{
			"participant_id":          "participant.code",
			"participant_name":        "customer.first_name",
			"company_name":            "company.first_name",
			"investment_product_name": "investment_product.name",
			"nab":                     "a.nab",
			"fee":                     "a.fee",
			"investment_at":           "a.investment_at",
		},
	})

	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]entity.ReportTransactionChannel)
	var data []*entity.ReportTransactionChannel = make([]*entity.ReportTransactionChannel, 0)

	// Copy data dan lakukan post-processing jika diperlukan
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
