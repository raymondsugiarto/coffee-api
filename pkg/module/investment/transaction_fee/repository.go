package transactionfee

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.TransactionFeeDto) (*entity.TransactionFeeDto, error)
	GetOperationalFeesByCompanyAndDateRange(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.TransactionFeeDto) (*entity.TransactionFeeDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.TransactionFeeDto).FromModel(m), nil
}

func (r *repository) GetOperationalFeesByCompanyAndDateRange(ctx context.Context, companyID string, startDate, endDate time.Time) (map[string]float64, error) {
	var results []struct {
		InvestmentProductID string
		TotalOperationalFee float64
	}

	err := r.db.WithContext(ctx).
		Model(&entity.TransactionFeeDto{}).
		Select("investment_product_id, COALESCE(SUM(operation_fee), 0) as total_operational_fee").
		Where("company_id = ?", companyID).
		Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).
		Where("deleted_at IS NULL").
		Group("investment_product_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	operationalFees := make(map[string]float64)
	for _, result := range results {
		operationalFees[result.InvestmentProductID] = result.TotalOperationalFee
	}

	return operationalFees, nil
}
