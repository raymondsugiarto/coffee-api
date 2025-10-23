package permission

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	FindAll(ctx context.Context, req *entity.PermissionFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) FindAll(ctx context.Context, req *entity.PermissionFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Permission = make([]model.Permission, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Permission{})
		if req.RoleID != "" {
			if req.AllPermission {
				q.Preload("RolePermissions")
				q.Joins("LEFT JOIN role_permission ON role_permission.permission_id = permission.id AND role_permission.deleted_at IS NULL AND role_permission.role_id = ?", req.RoleID)
			} else {
				q.Joins("JOIN role_permission ON role_permission.permission_id = permission.id").
					Where("role_permission.role_id = ?", req.RoleID)
			}
		} else if req.UserID != "" {
			q.Joins("JOIN role_permission ON role_permission.permission_id = permission.id AND role_permission.deleted_at IS NULL").
				Joins("JOIN role ON role.id = role_permission.role_id").
				Joins("JOIN user_has_role ON user_has_role.role_id = role.id").
				Where("user_has_role.user_id = ?", req.UserID)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Permission)
	var data []*entity.PermissionDto = make([]*entity.PermissionDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.PermissionDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
