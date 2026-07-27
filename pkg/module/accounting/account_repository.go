package accounting

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

// accountRepository is the chart-of-accounts CRUD layer. Mirrors
// the `salarycomponent` shape: simple PK lookups + an org-scoped
// FindAll that lets the upstream listing filter by name and code.
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository builds a Repository bound to the given DB.
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewAccountDtoFromModel(m), nil
}

func (r *accountRepository) Get(ctx context.Context, id string) (*entity.AccountDto, error) {
	var m model.Account
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewAccountDtoFromModel(&m), nil
}

func (r *accountRepository) Update(ctx context.Context, dto *entity.AccountDto) (*entity.AccountDto, error) {
	if err := r.db.WithContext(ctx).Save(dto.ToModel()).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	// We rely on the (organization_id, code) UNIQUE constraint and
	// the FK from account_mutation.account_id to surface a clean
	// error if the row is in use. The repository intentionally does
	// NOT pre-check refcount — doing so would race with new
	// mutations inserting concurrently. Let the DB enforce.
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Account{}).Error
}

func (r *accountRepository) FindAll(
	ctx context.Context,
	req *entity.AccountFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.Account = make([]model.Account, 0)
	tbl := pagination.NewTable(r.db.WithContext(ctx))
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Account{})
		// NULL organization_id rows are global (visible to every
		// org) — same convention as item_category, item, and
		// salary_component.
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
		AllowedFields: []string{"code", "name"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.Account)
	out := make([]*entity.AccountDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewAccountDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
