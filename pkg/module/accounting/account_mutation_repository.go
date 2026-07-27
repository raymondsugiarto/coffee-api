package accounting

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

// accountMutationRepository is the ledger side of accounting.
//
// We deliberately keep this repository small: Create for upstream
// flows to post to the ledger, FindAll for reports to read it back.
// No Update / Delete because the ledger is append-only — fixing a
// posting is done by posting a counter-mutation, never by editing
// the original row.
type accountMutationRepository struct {
	db *gorm.DB
}

func NewAccountMutationRepository(db *gorm.DB) AccountMutationRepository {
	return &accountMutationRepository{db: db}
}

func (r *accountMutationRepository) Create(
	ctx context.Context,
	dto *entity.AccountMutationDto,
) (*entity.AccountMutationDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return entity.NewAccountMutationDtoFromModel(m), nil
}

func (r *accountMutationRepository) FindAll(
	ctx context.Context,
	req *entity.AccountMutationFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.AccountMutation = make([]model.AccountMutation, 0)
	tbl := pagination.NewTable(r.db.WithContext(ctx))
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.AccountMutation{})
		// Org-scoped like the rest of the catalog. NULL
		// organization_id rows are the seed / global view.
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where(
				"organization_id IS NULL OR organization_id = ?",
				req.FindAllRequest.OrganizationData.ID,
			)
		}
		return q
	}, &pagination.TableRequest{
		Request:    req,
		QueryField: []string{},
		Data:       &rows,
		AllowedFields: []string{
			"account_id", "ref_id", "ref_table", "ref_module", "amount", "created_at",
		},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.AccountMutation)
	out := make([]*entity.AccountMutationDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewAccountMutationDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
