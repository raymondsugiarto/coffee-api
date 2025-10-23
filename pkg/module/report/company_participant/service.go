package companyparticipant

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	GetCompanyParticipantReport(ctx context.Context, filter *entity.CompanyParticipantFilter) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCompanyParticipantReport(ctx context.Context, filter *entity.CompanyParticipantFilter) (*pagination.ResultPagination, error) {
	return s.repo.GetCompanyParticipantReport(ctx, filter)
}
