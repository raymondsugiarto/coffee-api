package reward

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error)
	Get(ctx context.Context, id string) (*entity.RewardDto, error)
	GetWithLock(ctx context.Context, tx *gorm.DB, id string) (*entity.RewardDto, error)
	Update(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RewardDto) (*entity.RewardDto, error)
	DecrementStock(ctx context.Context, tx *gorm.DB, id string) error
	IncrementStock(ctx context.Context, tx *gorm.DB, id string) error
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RewardDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.RewardDto, error) {
	var m model.Reward
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RewardDto).FromModel(&m), nil
}

func (r *repository) GetWithLock(ctx context.Context, tx *gorm.DB, id string) (*entity.RewardDto, error) {
	var m model.Reward
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RewardDto).FromModel(&m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error) {
	m := dto.ToModel()
	err := r.db.Model(&model.Reward{}).Where("id = ?", dto.ID).Updates(m).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *entity.RewardDto) (*entity.RewardDto, error) {
	m := dto.ToModel()
	err := tx.Model(&model.Reward{}).Where("id = ?", dto.ID).Updates(m).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) DecrementStock(ctx context.Context, tx *gorm.DB, id string) error {
	result := tx.Model(&model.Reward{}).
		Where("id = ? AND stock > 0", id).
		Update("stock", gorm.Expr("stock - 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) IncrementStock(ctx context.Context, tx *gorm.DB, id string) error {
	result := tx.Model(&model.Reward{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Reward{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Reward = make([]model.Reward, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Reward{}).
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		Data:          &m,
		AllowedFields: []string{"name", "points", "stock", "created_at"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Reward)
	var data []*entity.RewardDto = make([]*entity.RewardDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.RewardDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
