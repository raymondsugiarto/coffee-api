package portfolio

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllPortfolioWithNav(ctx context.Context, req *entity.PortfolioFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAllPortfolioWithNav(ctx context.Context, req *entity.PortfolioFindAllRequest) (*pagination.ResultPagination, error) {
	var m []*entity.PortfolioReportDto = make([]*entity.PortfolioReportDto, 0)

	date := time.Now().UTC()
	today := date.Format("2006-01-02")

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Table("unit_link").
			Select(`
				CONCAT(unit_link.participant_id, '-', unit_link.investment_product_id) as id,
				unit_link.participant_id,
				MAX(unit_link.customer_id) as customer_id,
				unit_link.investment_product_id,
				SUM(unit_link.ip) as ip,
				COALESCE(MAX(nav.amount), 0) as latest_nav,
				(SUM(unit_link.ip) * COALESCE(MAX(nav.amount), 0)) as total_balance,
				MAX(customer.first_name) as customer_first_name,
				MAX(participant.code) as participant_code,
				MAX(investment_product.code) as investment_product_code,
				MAX(investment_product.name) as investment_product_name
			`).
			Joins("LEFT JOIN customer ON customer.id = unit_link.customer_id").
			Joins("LEFT JOIN participant ON participant.id = unit_link.participant_id").
			Joins("LEFT JOIN investment_product ON investment_product.id = unit_link.investment_product_id").
			Joins("LEFT JOIN LATERAL (SELECT amount FROM net_asset_value WHERE investment_product_id = unit_link.investment_product_id AND deleted_at IS NULL AND created_date <= ? ORDER BY created_date DESC LIMIT 1) nav ON true", today).
			Where("unit_link.deleted_at is null").
			Group("unit_link.participant_id, unit_link.investment_product_id")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"participant_id", "customer_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	portfolioResults := result.Data.(*[]*entity.PortfolioReportDto)

	// Populate nested objects for each portfolio result
	for _, portfolio := range *portfolioResults {
		portfolio.PopulateNestedObjects()
	}

	// Update the result data with populated objects
	result.Data = *portfolioResults

	return result, nil
}
