package netassetvalue

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error)
	CreateBatchWithCallback(ctx context.Context, dto *entity.NetAssetValueBatchDto, cb func(tx *gorm.DB) error) (*entity.NetAssetValueBatchDto, error)
	FindByInvestmentProductAndDate(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error)
	FindByDate(ctx context.Context, date string) ([]*entity.NetAssetValueDto, error)
	Get(ctx context.Context, id string) (*entity.NetAssetValueDto, error)
	Update(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error)
	UpdateWithTx(ctx context.Context, dto *entity.NetAssetValueDto, tx *gorm.DB) (*entity.NetAssetValueDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.NetAssetValueFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateBatchWithCallback(ctx context.Context, dto *entity.NetAssetValueBatchDto, cb func(tx *gorm.DB) error) (*entity.NetAssetValueBatchDto, error) {
	m := dto.ToModel()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.CreateInBatches(m, len(dto.Items)).Error
		if err != nil {
			return err
		}

		return cb(tx)
	})
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (r *repository) FindByDate(ctx context.Context, date string) ([]*entity.NetAssetValueDto, error) {
	var m []*model.NetAssetValue

	netAssetValueDtos := make([]*entity.NetAssetValueDto, 0)
	if err := r.db.
		Where("created_date = ?", date).
		Where("deleted_at IS NULL").
		Find(&m).Error; err != nil {
		return netAssetValueDtos, err
	}
	for _, item := range m {
		netAssetValueDtos = append(netAssetValueDtos, new(entity.NetAssetValueDto).FromModel(item))
	}
	return netAssetValueDtos, nil
}

func (r *repository) FindByInvestmentProductAndDate(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error) {
	var m *model.NetAssetValue
	err := r.db.Where("investment_product_id = ? AND created_date = ?", investmentProductID, date).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.NetAssetValueDto).FromModel(m), nil
}

func (r *repository) Create(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.NetAssetValueDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.NetAssetValueDto, error) {
	var m *model.NetAssetValue
	err := r.db.Where("id = ?", id).Preload("InvestmentProduct").Preload("Customer").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.NetAssetValueDto).FromModel(m), nil
}

func (r *repository) UpdateWithTx(ctx context.Context, dto *entity.NetAssetValueDto, tx *gorm.DB) (*entity.NetAssetValueDto, error) {
	err := tx.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Update(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.NetAssetValue{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.NetAssetValueFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.NetAssetValue = make([]model.NetAssetValue, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.NetAssetValue{}).
			Preload("InvestmentProduct")
		if req.FormCreate {
			q.Joins("RIGHT JOIN investment_product ON investment_product.id = net_asset_value.investment_product_id AND created_date >= ?", time.Now().AddDate(0, 0, -1))
			q.Select("COALESCE(net_asset_value.investment_product_id, investment_product.id) as investment_product_id, net_asset_value.amount, net_asset_value.created_date")
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"created_date"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.NetAssetValue)
	var data []*entity.NetAssetValueDto = make([]*entity.NetAssetValueDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.NetAssetValueDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
