package transactionhistory

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionHistory(ctx context.Context, filter *entity.TransactionHistoryFilter) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTransactionHistory(ctx context.Context, req *entity.TransactionHistoryFilter) (*pagination.ResultPagination, error) {
	var m []entity.TransactionHistoryReport = make([]entity.TransactionHistoryReport, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.WithContext(ctx).Table("investment_item").
			Select(`
				COALESCE(investment_item.code, investment_item.id) as investment_item_code,
				investment_item.amount,
				investment_item.fee_amount,
				investment_item.total_amount,
				investment_item.investment_at,
				investment.status, 
				participant.code as participant_code,
				CONCAT(customer.first_name, ' ', customer.last_name) as participant_name,
				company.company_code as company_code,
				company.first_name as company_name,
				investment_product.code as investment_product_code,
				investment_product.name as investment_product_name,
				net_asset_value.amount as nav_amount
			`).
			Joins("JOIN investment ON investment.id = investment_item.investment_id").
			Joins("LEFT JOIN participant ON participant.id = investment_item.participant_id").
			Joins("LEFT JOIN customer ON customer.id = investment_item.customer_id").
			Joins("LEFT JOIN company ON company.id = customer.company_id").
			Joins("LEFT JOIN investment_product ON investment_product.id = investment_item.investment_product_id").
			Joins("LEFT JOIN net_asset_value ON investment_item.investment_product_id = net_asset_value.investment_product_id AND net_asset_value.created_date = date(investment_item.investment_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta')").
			Where("investment_item.deleted_at IS NULL")

		return q
	}, &pagination.TableRequest{
		Request: req,
		QueryField: []string{
			"participant.code",
			"company.company_code",
			"customer.first_name",
			"company.first_name",
		},
		Data:          &m,
		AllowedFields: []string{"investment_at", "status", "amount", "company_id"},
		MapFields: map[string]string{
			"investment_at": "investment_item.investment_at",
			"status":        "investment.status",
			"amount":        "investment_item.amount",
			"company_id":    "company.id",
		},
	})

	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]entity.TransactionHistoryReport)
	var data []*entity.TransactionHistoryReport = make([]*entity.TransactionHistoryReport, 0)
	for _, item := range *results {
		itemCopy := item
		if item.NavAmount != nil && *item.NavAmount > 0 {
			itemCopy.UnitAmount = item.TotalAmount / *item.NavAmount
		} else {
			itemCopy.UnitAmount = 0
		}
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
