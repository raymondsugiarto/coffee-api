package benefit_type

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	Create(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error)
	FindByID(ctx context.Context, id string) (*entity.BenefitTypeDto, error)
	FindAll(ctx context.Context, req *entity.BenefitTypeFindAllRequest) (*pagination.ResultPagination, error)
	Update(ctx context.Context, id string, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error)
	Delete(ctx context.Context, id string) error
	GetActiveBenefitTypeNames(ctx context.Context) ([]string, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error) {
	created, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.BenefitTypeDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.BenefitTypeFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) Update(ctx context.Context, id string, dto *entity.BenefitTypeDto) (*entity.BenefitTypeDto, error) {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dto.ID = id

	return s.repository.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return errors.New("benefit type not found")
	}

	return s.repository.Delete(ctx, id)
}

func (s *service) GetActiveBenefitTypeNames(ctx context.Context) ([]string, error) {
	req := &entity.BenefitTypeFindAllRequest{
		Status: "ACTIVE",
	}

	result, err := s.repository.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, item := range result.Data.([]*entity.BenefitTypeDto) {
		names = append(names, item.Name)
	}

	return names, nil
}
