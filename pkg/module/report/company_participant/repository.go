package companyparticipant

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	GetCompanyParticipantReport(ctx context.Context, filter *entity.CompanyParticipantFilter) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCompanyParticipantReport(ctx context.Context, req *entity.CompanyParticipantFilter) (*pagination.ResultPagination, error) {
	var m []entity.CompanyParticipantReport = make([]entity.CompanyParticipantReport, 0)

	contributionSubQuery := r.db.Table("investment_item").
		Select("investment_item.participant_id, SUM(investment_item.amount) as total_contribution").
		Joins("JOIN investment ON investment.id = investment_item.investment_id").
		Where("investment.status = 'SUCCESS'").
		Where("investment_item.deleted_at IS NULL").
		Group("investment_item.participant_id")

	unitSubQuery := r.db.Table("unit_link").
		Select("participant_id, SUM(ip) as total_unit").
		Where("deleted_at IS NULL").
		Group("participant_id")

	latestNavSubQuery := r.db.Table("net_asset_value nav1").
		Select("nav1.investment_product_id, nav1.amount").
		Where("nav1.deleted_at IS NULL").
		Where("nav1.created_date = (SELECT MAX(nav2.created_date) FROM net_asset_value nav2 WHERE nav2.investment_product_id = nav1.investment_product_id AND nav2.deleted_at IS NULL)")

	balanceSubQuery := r.db.Table("unit_link").
		Select("unit_link.participant_id, SUM(unit_link.ip * COALESCE(latest_nav.amount, 0)) as last_balance").
		Joins("LEFT JOIN (?) latest_nav ON latest_nav.investment_product_id = unit_link.investment_product_id", latestNavSubQuery).
		Where("unit_link.deleted_at IS NULL").
		Group("unit_link.participant_id")

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.WithContext(ctx).Table("participant").
			Select(`
				participant.id as participant_id,
				customer.id as customer_id,
				CONCAT(customer.first_name, ' ', COALESCE(customer.last_name, '')) as customer_name,
				participant.code as participant_code,
				COALESCE(contribution.total_contribution, 0) as total_contribution,
				COALESCE(unit.total_unit, 0) as total_unit,
				COALESCE(balance.last_balance, 0) as last_balance
			`).
			Joins("JOIN customer ON customer.id = participant.customer_id").
			Joins("JOIN company ON company.id = customer.company_id").
			Joins("LEFT JOIN (?) contribution ON contribution.participant_id = participant.id", contributionSubQuery).
			Joins("LEFT JOIN (?) unit ON unit.participant_id = participant.id", unitSubQuery).
			Joins("LEFT JOIN (?) balance ON balance.participant_id = participant.id", balanceSubQuery).
			Where("participant.deleted_at IS NULL").
			Where("customer.deleted_at IS NULL")

		return q
	}, &pagination.TableRequest{
		Request:    req,
		QueryField: []string{},
		Data:       &m,
		AllowedFields: []string{
			"company.id",
		},
		MapFields: map[string]string{
			"created_at": "participant.created_at",
		},
	})

	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]entity.CompanyParticipantReport)
	var data []*entity.CompanyParticipantReport = make([]*entity.CompanyParticipantReport, 0)
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
