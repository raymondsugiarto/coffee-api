package cashdebt

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error)
	Get(ctx context.Context, id string) (*entity.CashDebtDto, error)
	Update(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.CashDebtFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewCashDebtDtoFromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.CashDebtDto, error) {
	var m model.CashDebt
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewCashDebtDtoFromModel(&m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.CashDebtDto) (*entity.CashDebtDto, error) {
	if err := r.db.WithContext(ctx).Save(dto.ToModel()).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.CashDebt{}).Error
}

func (r *repository) FindAll(
	ctx context.Context,
	req *entity.CashDebtFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.CashDebt = make([]model.CashDebt, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.CashDebt{})
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &rows,
		AllowedFields: []string{"date", "amount", "payment_method", "created_at"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.CashDebt)
	out := make([]*entity.CashDebtDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewCashDebtDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
