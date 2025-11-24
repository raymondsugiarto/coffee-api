package order

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.OrderDto) (*entity.OrderDto, error)
	FindByID(ctx context.Context, id string) (*entity.OrderDto, error)
	Update(ctx context.Context, dto *entity.OrderDto) (*entity.OrderDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.OrderFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo           Repository
	companyService company.Service
}

func NewService(repo Repository, companyService company.Service) Service {
	return &service{repo, companyService}
}

func (s *service) Create(ctx context.Context, dto *entity.OrderDto) (*entity.OrderDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	company, err := s.companyService.FindCompanyByAdminID(ctx, dto.AdminID)
	if err != nil {
		return nil, err
	}
	dto.OrderAt = time.Now()
	code, _ := gonanoid.Generate("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 4)
	dto.Code = time.Now().Format("20060102") + "/" + code
	dto.CompanyID = company.ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.OrderDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.OrderDto) (*entity.OrderDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.OrderFindAllRequest) (*pagination.ResultPagination, error) {
	company, err := s.companyService.FindCompanyByAdminID(ctx, req.AdminID)
	if err != nil {
		return nil, err
	}
	req.CompanyID = company.ID
	return s.repo.FindAll(ctx, req)
}
