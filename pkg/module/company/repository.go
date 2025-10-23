package company

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.CompanyDto, cb func(tx *gorm.DB) error) (*entity.CompanyDto, error)
	FindByEmail(ctx context.Context, email string) (*entity.CompanyDto, error)
	FindByUserID(ctx context.Context, userId string) (*entity.CompanyDto, error)
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.CompanyDto, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, dto *entity.CompanyDto) (*entity.CompanyDto, error)
	FindByCompanyCode(ctx context.Context, companyCode string) (*entity.CompanyDto, error)
	CountByType(ctx context.Context) ([]*entity.CountCompanyPerType, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.CompanyDto, cb func(tx *gorm.DB) error) (*entity.CompanyDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	m := dto.ToModel()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		dto.ID = m.ID
		dto.UserID = m.UserID
		dto.User.ID = m.UserID
		if err := cb(tx); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*entity.CompanyDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Company
	if err := r.db.
		Preload("User").
		Preload("DomisiliObject").
		Where("email = ? and organization_id = ?", email, organizationID).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) FindByUserID(ctx context.Context, userId string) (*entity.CompanyDto, error) {
	var m *model.Company
	err := r.db.
		Preload("User").
		Preload("DomisiliObject").
		Where("user_id = ?", userId).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Company = make([]model.Company, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.Company{}).
			Preload("User").
			Preload("DomisiliObject").
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{"first_name", "address", "pic_name", "pic_email", "pic_phone"},
		Data:          &m,
		AllowedFields: []string{"first_name", "last_name", "company_code"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Company)
	var data []*entity.CompanyDto = make([]*entity.CompanyDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.CompanyDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.CompanyDto, error) {
	var m *model.Company
	if err := r.db.
		Preload("User").
		Preload("DomisiliObject").
		Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Company{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, dto *entity.CompanyDto) (*entity.CompanyDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	m := dto.ToModel()
	err := r.db.Updates(m).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) FindByCompanyCode(ctx context.Context, companyCode string) (*entity.CompanyDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var m *model.Company
	if err := r.db.
		Preload("User").
		Preload("DomisiliObject").
		Where("company_code = ? and organization_id = ?", companyCode, organizationID).First(&m).Error; err != nil {
		return nil, err
	}
	return new(entity.CompanyDto).FromModel(m), nil
}

func (r *repository) CountByType(ctx context.Context) ([]*entity.CountCompanyPerType, error) {
	var m []*entity.CountCompanyPerType

	query := r.db.WithContext(ctx).Model(&model.Company{}).
		Select(`company.company_type, COUNT(id) as count`).
		Group("company.company_type")

	err := query.Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
