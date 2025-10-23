package investmentitem

import (
	"context"
	"time"

	customer "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type NavWithDate struct {
	Amount float64   `json:"amount"`
	Date   time.Time `json:"date"`
}

type Repository interface {
	Create(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error)
	CreateBatchWithTx(ctx context.Context, db *gorm.DB, dto []*entity.InvestmentItemDto) ([]*entity.InvestmentItemDto, error)
	Get(ctx context.Context, id string) (*entity.InvestmentItemDto, error)
	GetByCode(ctx context.Context, code string) (*entity.InvestmentItemDto, error)
	Update(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentItemFindAllRequest) (*pagination.ResultPagination, error)
	SumInvestmentByCustomer(ctx context.Context, customerID string) (*customer.SumUnitLinkPortfolioDto, error)
	FindByInvestmentID(ctx context.Context, investmentID string) ([]*entity.InvestmentItemDto, error)
	HasPreviousInvestmentForCustomer(ctx context.Context, investmentID string, customerID string) (bool, error)
	GetLatestNavByProductIDs(ctx context.Context, productIDs []string) (map[string]float64, error)
	GetNavByProductIDsWithDate(ctx context.Context, productIDs []string, dateString string) (map[string]NavWithDate, error)
}

const (
	BATCH_SIZE = 100
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateBatchWithTx(ctx context.Context, db *gorm.DB, dto []*entity.InvestmentItemDto) ([]*entity.InvestmentItemDto, error) {
	var m []model.InvestmentItem = make([]model.InvestmentItem, 0)
	for _, item := range dto {
		m = append(m, *item.ToModel())
	}
	err := db.CreateInBatches(&m, BATCH_SIZE).Error
	if err != nil {
		return nil, err
	}
	var data []*entity.InvestmentItemDto = make([]*entity.InvestmentItemDto, 0)
	for _, item := range m {
		data = append(data, new(entity.InvestmentItemDto).FromModel(&item))
	}
	return data, nil
}

func (r *repository) Create(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentItemDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.InvestmentItemDto, error) {
	var m *model.InvestmentItem
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentItemDto).FromModel(m), nil
}

func (r *repository) GetByCode(ctx context.Context, code string) (*entity.InvestmentItemDto, error) {
	var m *model.InvestmentItem
	err := r.db.Where("code = ?", code).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentItemDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.InvestmentItemDto) (*entity.InvestmentItemDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.InvestmentItem{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.InvestmentItemFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.InvestmentItem = make([]model.InvestmentItem, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.InvestmentItem{}).
			Preload("Investment").
			Preload("InvestmentProduct").
			Preload("Participant").
			Preload("Customer")

		if req.InvestmentStatus != "" || req.CompanyID != "" {
			joinClause := "JOIN investment ON investment.id = investment_item.investment_id"
			joinArgs := []interface{}{}

			if req.InvestmentStatus != "" {
				joinClause += " AND investment.status = ?"
				joinArgs = append(joinArgs, req.InvestmentStatus)
			}

			q = q.Joins(joinClause, joinArgs...)

			if req.CompanyID != "" {
				q = q.Joins("JOIN company ON company.id = investment.company_id").
					Preload("Investment.Company").
					Where("company.id = ?", req.CompanyID)
			}
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"investment_at", "customer_id", "status", "investment_id"},
		MapFields: map[string]string{
			"investment_at": "investment_item.investment_at",
			"customer_id":   "investment_item.customer_id",
			"investment_id": "investment_id",
		},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.InvestmentItem)
	var data []*entity.InvestmentItemDto = make([]*entity.InvestmentItemDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.InvestmentItemDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) SumInvestmentByCustomer(ctx context.Context, customerID string) (*customer.SumUnitLinkPortfolioDto, error) {
	var m *customer.SumUnitLinkPortfolioDto

	subquery := r.db.WithContext(ctx).Model(&model.InvestmentItem{}).
		Joins("JOIN investment_payment ON investment_item.investment_id = investment_payment.investment_id").
		Where("investment_payment.status = ? AND investment_item.customer_id = ?", "success", customerID).
		Select("investment_item.amount")

	err := r.db.WithContext(ctx).Table("(?) as sub", subquery).
		Select("COALESCE(SUM(sub.amount), 0) as total_payment").
		Take(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) FindByInvestmentID(ctx context.Context, investmentID string) ([]*entity.InvestmentItemDto, error) {
	var m []model.InvestmentItem
	err := r.db.WithContext(ctx).
		Where("investment_id = ?", investmentID).
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var d []*entity.InvestmentItemDto
	for _, item := range m {
		d = append(d, new(entity.InvestmentItemDto).FromModel(&item))
	}
	return d, nil
}

func (r *repository) HasPreviousInvestmentForCustomer(ctx context.Context, investmentID string, customerID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.InvestmentItem{}).
		Where("investment_id != ? AND customer_id = ?", investmentID, customerID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) GetLatestNavByProductIDs(ctx context.Context, productIDs []string) (map[string]float64, error) {
	type NavResult struct {
		InvestmentProductID string  `gorm:"column:investment_product_id"`
		Amount              float64 `gorm:"column:amount"`
	}

	var results []NavResult
	subquery := r.db.WithContext(ctx).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Table("net_asset_value").
		Where("investment_product_id IN ?", productIDs)

	err := r.db.WithContext(ctx).
		Table("(?) as nav_ranked", subquery).
		Select("investment_product_id, amount").
		Where("rn = 1").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	navMap := make(map[string]float64)
	for _, result := range results {
		navMap[result.InvestmentProductID] = result.Amount
	}

	return navMap, nil
}

func (r *repository) GetNavByProductIDsWithDate(ctx context.Context, productIDs []string, targetDate string) (map[string]NavWithDate, error) {
	type NavResult struct {
		InvestmentProductID string    `gorm:"column:investment_product_id"`
		Amount              float64   `gorm:"column:amount"`
		CreatedDate         time.Time `gorm:"column:created_date"`
	}

	var results []NavResult

	subquery := r.db.WithContext(ctx).
		Select("investment_product_id, amount, created_date, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Table("net_asset_value").
		Where("investment_product_id IN ?", productIDs).
		Where("created_date <= ?", targetDate)

	err := r.db.WithContext(ctx).
		Table("(?) as nav_ranked", subquery).
		Select("investment_product_id, amount, created_date").
		Where("rn = 1").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	navMap := make(map[string]NavWithDate)
	for _, result := range results {
		navMap[result.InvestmentProductID] = NavWithDate{
			Amount: result.Amount,
			Date:   result.CreatedDate,
		}
	}

	return navMap, nil
}
