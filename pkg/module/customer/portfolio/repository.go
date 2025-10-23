package portfolio

import (
	"context"

	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error)
	Get(ctx context.Context, id string) (*entity.PortfolioDto, error)
	Update(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error) {
	m := dto.ToModel()
	err := r.db.WithContext(ctx).Where(model.Portfolio{
		ParticipantID:       m.ParticipantID,
		InvestmentProductID: m.InvestmentProductID,
	}).FirstOrCreate(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.PortfolioDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.PortfolioDto, error) {
	var m *model.Portfolio
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.PortfolioDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.PortfolioDto) (*entity.PortfolioDto, error) {
	err := r.db.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Portfolio{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *e.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Portfolio = make([]model.Portfolio, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Portfolio{}).
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
	results := result.Data.(*[]model.Portfolio)
	var data []*entity.PortfolioDto = make([]*entity.PortfolioDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.PortfolioDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
