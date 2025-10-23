package rolepermission

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type Service interface {
	AddPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) AddPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	dto := &entity.RolePermissionDto{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return s.repo.AddPermissionToRole(ctx, dto)
}

func (s *service) RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error {
	dto, err := s.repo.FindByRoleIDAndPermissionID(ctx, roleID, permissionID)
	if err != nil {
		return err
	}
	return s.repo.RemovePermissionFromRole(ctx, dto)
}
