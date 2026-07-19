package company

import (
	"context"
	"strings"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error)
	FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error)
	FindAll(ctx context.Context, req *entity.CompanyFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error) {
	var m *model.Company
	err := r.db.Model(model.Company{}).
		Joins("JOIN admin_company ON company.id = admin_company.company_id").
		Joins("JOIN admin ON admin.id = admin_company.admin_id").
		Where("admin.user_id = ?", userID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewCompanyDtoFromModel(m), nil
}

func (r *repository) FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error) {
	var m *model.Company
	err := r.db.Model(model.Company{}).
		Joins("JOIN admin_company ON company.id = admin_company.company_id").
		Where("admin_company.admin_id = ?", adminID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return entity.NewCompanyDtoFromModel(m), nil
}

// FindAll paginates the company table. Companies are org-scoped:
// NULL organization_id rows are treated as global seed data so
// every org sees the master list (e.g. SEKIAN, Mawaru) seeded in
// 000006_salary_component. The Query filter does a case-insensitive
// match on `name`.
func (r *repository) FindAll(
	ctx context.Context,
	req *entity.CompanyFindAllRequest,
) (*pagination.ResultPagination, error) {
	var rows []model.Company = make([]model.Company, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Company{})
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where(
				"organization_id IS NULL OR organization_id = ?",
				req.FindAllRequest.OrganizationData.ID,
			)
		}
		if s := strings.TrimSpace(req.Query); s != "" {
			like := "%" + s + "%"
			q = q.Where("name ILIKE ?", like)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &rows,
		AllowedFields: []string{"name"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	hits := result.Data.(*[]model.Company)
	out := make([]*entity.CompanyDto, 0, len(*hits))
	for i := range *hits {
		out = append(out, entity.NewCompanyDtoFromModel(&(*hits)[i]))
	}
	return &pagination.ResultPagination{
		Data:        out,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
