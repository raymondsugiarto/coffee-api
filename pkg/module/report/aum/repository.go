package aum

import (
	"context"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetCompanyAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]entity.ReportAUMData, error)
	GetPesertaMandiriAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]entity.ReportAUMData, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetCompanyAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]entity.ReportAUMData, error) {
	var result []entity.ReportAUMData

	subQuery, dateStr := r.getAUMSubQueryAndDate(ctx, filter)

	err := r.db.WithContext(ctx).
		Table("company").
		Select("company.first_name as company_name, company.pilar_type, COALESCE(SUM(unit_aum.aum), 0) as total_aum").
		Joins("LEFT JOIN customer ON customer.company_id = company.id AND customer.deleted_at IS NULL").
		Joins("LEFT JOIN participant ON participant.customer_id = customer.id AND participant.deleted_at IS NULL").
		Joins("LEFT JOIN (?) unit_aum ON unit_aum.participant_id = participant.id", subQuery).
		Where("company.deleted_at IS NULL").
		Where("company.status = ?", "APPROVED").
		Where("company.company_type = ?", filter.CompanyType).
		Where("company.created_at < ?", dateStr).
		Group("company.id, company.first_name, company.pilar_type").
		Order("CASE WHEN company.pilar_type = 'PILAR' THEN 1 ELSE 2 END, company.first_name").
		Scan(&result).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get company AUM: %w", err)
	}

	return result, nil
}

// GetPesertaMandiriAUM retrieves AUM data for individual participants (customers without company_id)
// TODO: Current implementation assigns all peserta mandiri to PILAR category
// Business logic needs to be implemented to properly split between PILAR and NON PILAR
// based on investment_product type or other business rules
func (r *repository) GetPesertaMandiriAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]entity.ReportAUMData, error) {
	var result []entity.ReportAUMData

	subQuery, _ := r.getAUMSubQueryAndDate(ctx, filter)

	var pilarData entity.ReportAUMData
	err := r.db.WithContext(ctx).
		Table("customer").
		Select("'PESERTA MANDIRI' as company_name, 'PILAR' as pilar_type, COALESCE(SUM(unit_aum.aum), 0) as total_aum").
		Joins("LEFT JOIN participant ON participant.customer_id = customer.id AND participant.deleted_at IS NULL").
		Joins("LEFT JOIN (?) unit_aum ON unit_aum.participant_id = participant.id", subQuery).
		Where("customer.deleted_at IS NULL").
		Where("customer.sim_status = ?", model.SIMStatusActive).
		Where("(customer.company_id IS NULL OR customer.company_id = '')").
		Scan(&pilarData).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get peserta mandiri PILAR AUM: %w", err)
	}

	if pilarData.TotalAUM > 0 {
		result = append(result, pilarData)
	}

	// Query for NON PILAR - currently set to 0 as per TODO
	// TODO: Implement proper business logic to split PILAR vs NON PILAR based on investment_product type
	// Note: This block is placeholder; TotalAUM is hardcoded to 0, so append never executes until TODO is addressed.
	var nonPilarData entity.ReportAUMData
	nonPilarData.CompanyName = "PESERTA MANDIRI"
	nonPilarData.PilarType = "NON PILAR"
	nonPilarData.TotalAUM = 0

	if nonPilarData.TotalAUM > 0 {
		result = append(result, nonPilarData)
	}

	return result, nil
}

// getAUMSubQueryAndDate calculates the start of next month and builds the AUM subquery
// for unit_link with filter before start of next month for end-of-period calculation.
func (r *repository) getAUMSubQueryAndDate(ctx context.Context, filter *entity.ReportAUMFilter) (*gorm.DB, string) {
	nextMonth := filter.Month + 1
	nextYear := filter.Year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	startOfNextMonth := time.Date(nextYear, time.Month(nextMonth), 1, 0, 0, 0, 0, time.UTC)
	startOfNextMonthStr := startOfNextMonth.Format("2006-01-02")

	subQuery := r.db.WithContext(ctx).
		Table("unit_link").
		Select("unit_link.participant_id, SUM(unit_link.ip * COALESCE(latest_nav.amount, 0)) AS aum").
		Joins(`LEFT JOIN (
			SELECT investment_product_id, amount,
				ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn
			FROM net_asset_value
			WHERE deleted_at IS NULL
			AND EXTRACT(MONTH FROM created_date) = ?
			AND EXTRACT(YEAR FROM created_date) = ?
		) latest_nav ON latest_nav.investment_product_id = unit_link.investment_product_id AND latest_nav.rn = 1`, filter.Month, filter.Year).
		Where("unit_link.deleted_at IS NULL").
		Where("unit_link.transaction_date < ?", startOfNextMonthStr).
		Group("unit_link.participant_id")

	return subQuery, startOfNextMonthStr
}
