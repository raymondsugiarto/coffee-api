package reward

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

type Service interface {
	Create(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error)
	FindByID(ctx context.Context, id string) (*entity.RewardDto, error)
	Update(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error) {
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.RewardDto, error) {
	reward, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return reward, nil
}

func (s *service) Update(ctx context.Context, dto *entity.RewardDto) (*entity.RewardDto, error) {
	_, err := s.repo.Get(ctx, dto.ID)
	if err != nil {
		return nil, status.New(status.EntityNotFound, errors.New("reward not found"))
	}

	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return status.New(status.EntityNotFound, errors.New("reward not found"))
	}

	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}
