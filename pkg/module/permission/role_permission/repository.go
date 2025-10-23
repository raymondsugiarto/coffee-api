package rolepermission

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	AddPermissionToRole(ctx context.Context, dto *entity.RolePermissionDto) error
	RemovePermissionFromRole(ctx context.Context, dto *entity.RolePermissionDto) error
	FindByRoleIDAndPermissionID(ctx context.Context, roleID string, permissionID string) (*entity.RolePermissionDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByRoleIDAndPermissionID(ctx context.Context, roleID string, permissionID string) (*entity.RolePermissionDto, error) {
	var m *model.RolePermission
	err := r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.RolePermissionDto).FromModel(m), nil
}

func (r *repository) AddPermissionToRole(ctx context.Context, dto *entity.RolePermissionDto) error {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RemovePermissionFromRole(ctx context.Context, dto *entity.RolePermissionDto) error {
	m := dto.ToModel()
	err := r.db.Delete(m).Error
	if err != nil {
		return err
	}
	return nil
}
