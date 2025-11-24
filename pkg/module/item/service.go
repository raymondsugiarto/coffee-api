package item

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error)
	FindByID(ctx context.Context, id string) (*entity.ItemDto, error)
	Update(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo           Repository
	companyService company.Service
}

func NewService(repo Repository, companyService company.Service) Service {
	return &service{repo, companyService}
}

func (s *service) Create(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.ItemDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.ItemDto) (*entity.ItemDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.ItemFindAllRequest) (*pagination.ResultPagination, error) {
	if req.MyEmployeeItem {
		company, err := s.companyService.FindCompanyByUserID(ctx, req.UserID)
		if err != nil {
			return nil, err
		}
		log.WithContext(ctx).Infof("company %+v", company)
		req.CompanyID = company.ID
	}
	return s.repo.FindAll(ctx, req)
}
