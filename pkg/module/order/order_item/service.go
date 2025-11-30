package orderitem

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
)

type Service interface {
	Count(ctx context.Context, req *entity.OrderFindAllRequest) ([]entity.OrderItemPerItemCountDto, error)
}

type service struct {
	repo           Repository
	companyService company.Service
}

func NewService(repo Repository, companyService company.Service) Service {
	return &service{repo, companyService}
}

func (s *service) Count(ctx context.Context, req *entity.OrderFindAllRequest) ([]entity.OrderItemPerItemCountDto, error) {
	company, err := s.companyService.FindCompanyByAdminID(ctx, req.AdminID)
	if err != nil {
		log.WithContext(ctx).Errorf("error count %v", err)
		return nil, err
	}
	req.CompanyID = company.ID
	return s.repo.Count(ctx, req)
}
