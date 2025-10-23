package customerpoint

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error)
	CreateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error)
	Get(ctx context.Context, id string) (*entity.CustomerPointDto, error)
	Update(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	GetTotalPoint(ctx context.Context, customerID string) (float64, error)
	GetTotalPointWithLock(ctx context.Context, tx *gorm.DB, customerID string) (float64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.CustomerPointDto).FromModel(m), nil
}

func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error) {
	m := dto.ToModel()
	err := tx.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.CustomerPointDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.CustomerPointDto, error) {
	var m *model.CustomerPoint
	err := r.db.Where("id = ?", id).Preload("RoleParent").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.CustomerPointDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error) {
	err := r.db.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.CustomerPoint{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.CustomerPoint = make([]model.CustomerPoint, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.CustomerPoint{}).
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
	results := result.Data.(*[]model.CustomerPoint)
	var data []*entity.CustomerPointDto = make([]*entity.CustomerPointDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.CustomerPointDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) GetTotalPoint(ctx context.Context, customerID string) (float64, error) {
	var totalPoint float64
	err := r.db.Model(&model.CustomerPoint{}).
		Select("sum(point)").
		Where("customer_id = ?", customerID).
		Where("deleted_at is null").
		Group("customer_id").
		Scan(&totalPoint).Error
	if err != nil {
		return 0, err
	}
	return totalPoint, nil
}

func (r *repository) GetTotalPointWithLock(ctx context.Context, tx *gorm.DB, customerID string) (float64, error) {
	var customerPoints []model.CustomerPoint
	err := tx.Where("customer_id = ? AND deleted_at IS NULL", customerID).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&customerPoints).Error
	if err != nil {
		return 0, err
	}

	var totalPoint float64
	for _, point := range customerPoints {
		totalPoint += point.Point
	}
	return totalPoint, nil
}
