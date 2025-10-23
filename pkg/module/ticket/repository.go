package ticket

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error)
	CreateWithTx(ctx context.Context, dto *entity.TicketDto, cb func(tx *gorm.DB) error) (*entity.TicketDto, error)
	Get(ctx context.Context, id string) (*entity.TicketDto, error)
	Update(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.TicketDto).FromModel(m), nil
}

func (r *repository) CreateWithTx(ctx context.Context, dto *entity.TicketDto, cb func(tx *gorm.DB) error) (*entity.TicketDto, error) {
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
	return new(entity.TicketDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.TicketDto, error) {
	var m *model.Ticket
	err := r.db.Model(&model.Ticket{}).
		Preload("User").
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.TicketDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.TicketDto) (*entity.TicketDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Ticket{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Ticket = make([]model.Ticket, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Ticket{}).Preload("User").
			Joins("LEFT JOIN customer ON customer.user_id = ticket.user_id")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"customer.first_name", "customer.last_name"},
		Data:          &m,
		AllowedFields: []string{"customer.first_name", "customer.last_name"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	return &pagination.ResultPagination{
		Data:        &m,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
