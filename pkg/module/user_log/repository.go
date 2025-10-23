package userlog

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error)
	CreateWithTx(ctx context.Context, dto *entity.UserLogDto, cb func(tx *gorm.DB) error) (*entity.UserLogDto, error)
	Get(ctx context.Context, id string) (*entity.UserLogDto, error)
	Update(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.UserLogFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.UserLogDto).FromModel(m), nil
}

func (r *repository) CreateWithTx(ctx context.Context, dto *entity.UserLogDto, cb func(tx *gorm.DB) error) (*entity.UserLogDto, error) {
	m := dto.ToModel()
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		if err := cb(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return new(entity.UserLogDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.UserLogDto, error) {
	var m *model.UserLog
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.UserLogDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.UserLogDto) (*entity.UserLogDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.UserLog{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.UserLogFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.UserLog = make([]model.UserLog, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.UserLog{}).Preload("UserCredential").
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"company_id", "ref_module", "ref_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.UserLog)
	var data []*entity.UserLogDto = make([]*entity.UserLogDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.UserLogDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
