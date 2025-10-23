package bankcustomer

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.BankCustomerDto, error)
	Update(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error)
	SetAllNonDefault(ctx context.Context, customerID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.BankCustomerDto).FromModel(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.BankCustomerDto, error) {
	var m *model.BankCustomer
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.BankCustomerDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.BankCustomerDto) (*entity.BankCustomerDto, error) {
	err := r.db.Model(&model.BankCustomer{}).Where("id = ?", dto.ID).Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.BankCustomer{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.BankCustomerFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.BankCustomer = make([]model.BankCustomer, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.BankCustomer{}).
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"customer_id", "is_default"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.BankCustomer)
	var data []*entity.BankCustomerDto = make([]*entity.BankCustomerDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.BankCustomerDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) SetAllNonDefault(ctx context.Context, customerID string) error {
	err := r.db.Model(&model.BankCustomer{}).
		Where("customer_id = ?", customerID).
		Update("is_default", false).Error
	if err != nil {
		return err
	}
	return nil
}
