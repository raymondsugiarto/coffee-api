package salarycomponent

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error)
	Get(ctx context.Context, id string) (*entity.SalaryComponentDto, error)
	Update(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.SalaryComponentFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error) {
	m := dto.ToModel()
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewSalaryComponentDtoFromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.SalaryComponentDto, error) {
	var m model.SalaryComponent
	if err := r.db.Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewSalaryComponentDtoFromModel(&m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.SalaryComponentDto) (*entity.SalaryComponentDto, error) {
	if err := r.db.Save(dto.ToModel()).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SalaryComponent{}).Error
}

func (r *repository) FindAll(
	ctx context.Context,
	req *entity.SalaryComponentFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.SalaryComponent = make([]model.SalaryComponent, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.SalaryComponent{})
		// Org-scoped like stock_session + item + item_category.
		// NULL organization_id rows are treated as global seed.
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where(
				"organization_id IS NULL OR organization_id = ?",
				req.FindAllRequest.OrganizationData.ID,
			)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &rows,
		AllowedFields: []string{"company_id", "component_type", "minimum_target", "amount"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.SalaryComponent)
	out := make([]*entity.SalaryComponentDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewSalaryComponentDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
