package investmentproduct

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error)
	Get(ctx context.Context, id string) (*entity.InvestmentProductDto, error)
	Update(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentProductFilter) (*pagination.ResultPagination, error)

	SummaryList(ctx context.Context) ([]*entity.InvestmentProductSummaryDto, error)
	CalculateAUM(ctx context.Context, investmentProductID string) (float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentProductDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.InvestmentProductDto, error) {
	var m *model.InvestmentProduct
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentProductDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.InvestmentProductDto) (*entity.InvestmentProductDto, error) {
	err := r.db.Model(&model.InvestmentProduct{}).
		Where("id = ?", dto.ID).
		Select("*").
		Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.InvestmentProduct{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.InvestmentProductFilter) (*pagination.ResultPagination, error) {
	var m []model.InvestmentProduct = make([]model.InvestmentProduct, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.InvestmentProduct{}).
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{""},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.InvestmentProduct)
	var data []*entity.InvestmentProductDto = make([]*entity.InvestmentProductDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.InvestmentProductDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) SummaryList(ctx context.Context) ([]*entity.InvestmentProductSummaryDto, error) {
	var m []*entity.InvestmentProductSummaryDto

	date := time.Now().UTC()
	today := date.Format("2006-01-02")

	err := r.db.WithContext(ctx).Model(&model.InvestmentProduct{}).
		Joins("LEFT JOIN unit_link ON unit_link.investment_product_id = investment_product.id").
		Joins("LEFT JOIN net_asset_value ON net_asset_value.investment_product_id = investment_product.id AND net_asset_value.created_date = ?", today).
		Select("investment_product.id, investment_product.name as investment_product_name, SUM(ip) as sum_ip, SUM(ip * nab) as invest_amount, (SUM(ip) * MAX(net_asset_value.amount)) as total_amount, MAX(net_asset_value.amount) as nab").
		Group("investment_product.id, investment_product.name").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) CalculateAUM(ctx context.Context, investmentProductID string) (float64, error) {
	var out struct {
		AUM float64
	}

	valuationDate := time.Now().UTC().Format("2006-01-02")

	// Calculate total units using ip if available, fallback to total_amount / NAV on transaction date
	unitWithCalculatedIP := r.db.WithContext(ctx).
		Select(`ul.investment_product_id,
			CASE
				WHEN ul.ip > 0 THEN ul.ip
				WHEN nav_tx.amount > 0 THEN ul.total_amount / nav_tx.amount
				ELSE 0
			END AS calculated_ip`).
		Table("unit_link ul").
		Joins("LEFT JOIN net_asset_value nav_tx ON nav_tx.investment_product_id = ul.investment_product_id AND nav_tx.created_date = ul.transaction_date").
		Where("ul.investment_product_id = ? AND ul.deleted_at IS NULL AND ul.transaction_date <= ?", investmentProductID, valuationDate)

	unitSummary := r.db.WithContext(ctx).
		Select("investment_product_id, SUM(calculated_ip) AS total_ip").
		Table("(?) AS ucip", unitWithCalculatedIP).
		Group("investment_product_id")

	// Get latest NAV (created_date <= valuationDate) for the product
	latestNav := r.db.WithContext(ctx).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) AS rn").
		Table("net_asset_value").
		Where("investment_product_id = ? AND created_date <= ?", investmentProductID, valuationDate)

	// AUM = total_ip × latest NAV
	if err := r.db.WithContext(ctx).
		Table("(?) AS us", unitSummary).
		Joins("LEFT JOIN (?) AS ln ON ln.investment_product_id = us.investment_product_id AND ln.rn = 1", latestNav).
		Select("COALESCE(us.total_ip * COALESCE(ln.amount, 0), 0) AS aum").
		Take(&out).Error; err != nil {
		return 0, err
	}

	return out.AUM, nil
}
