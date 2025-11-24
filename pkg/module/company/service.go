package company

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
)

type Service interface {
	FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error)
	FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) FindCompanyByUserID(ctx context.Context, userID string) (*entity.CompanyDto, error) {
	return s.repository.FindCompanyByUserID(ctx, userID)
}

func (s *service) FindCompanyByAdminID(ctx context.Context, adminID string) (*entity.CompanyDto, error) {
	return s.repository.FindCompanyByAdminID(ctx, adminID)
}
