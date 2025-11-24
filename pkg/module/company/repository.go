package company

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error)
	FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error)
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
