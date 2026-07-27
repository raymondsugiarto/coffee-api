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
	// FindByCompany returns every salary_component row scoped to
	// the given company, ordered by minimum_target ASC. Used by the
	// stock-session close path to resolve the per-session salary
	// breakdown; the ascending order matches the picking logic in
	// entity.StockSessionDto.RecomputeSalary (the highest
	// minimum_target the driver still clears wins, so we want to
	// walk the list from low to high).
	FindByCompany(ctx context.Context, companyID string) ([]*entity.SalaryComponentDto, error)
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

// findByCompany returns every salary_component row for a company,
// ordered by minimum_target ASC then component_type ASC for a
// stable tie-break.
//
// Why ordered:
//   - The stock-session close path picks the highest
//     minimum_target that the driver's TotalItems still clears.
//     Walking the list ascending keeps the picking loop in O(n)
//     with no sort, and makes the row ordering deterministic for
//     fixtures + tests.
//   - reports that aggregate by component_type (meal vs attendance
//     vs bonus) get the same row order across calls.
//
// Soft-deleted rows (deleted_at IS NOT NULL) are skipped via the
// GORM default that the rest of the codebase relies on.
func (r *repository) FindByCompany(
	ctx context.Context,
	companyID string,
) ([]*entity.SalaryComponentDto, error) {
	if companyID == "" {
		return []*entity.SalaryComponentDto{}, nil
	}
	var rows []model.SalaryComponent
	if err := r.db.WithContext(ctx).
		Where("company_id = ?", companyID).
		Order("minimum_target ASC, component_type ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]*entity.SalaryComponentDto, 0, len(rows))
	for i := range rows {
		out = append(out, entity.NewSalaryComponentDtoFromModel(&rows[i]))
	}
	return out, nil
}
