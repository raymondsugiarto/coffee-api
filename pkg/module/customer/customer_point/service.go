package customerpoint

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

type Service interface {
	Create(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error)
	Get(ctx context.Context, id string) (*entity.CustomerPointDto, error)
	Update(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	GetTotalPoint(ctx context.Context, customerID string) (float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error) {
	// TODO: need locking to prevent
	totalPoint, err := s.repo.GetTotalPoint(ctx, dto.CustomerID)
	if err != nil {
		return nil, err
	}
	if totalPoint+dto.Point < 0 {
		return nil, status.New(status.BadRequest, errors.New("total point cannot be less than 0"))
	}

	m, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}
func (s *service) Get(ctx context.Context, id string) (*entity.CustomerPointDto, error) {
	m, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) Update(ctx context.Context, dto *entity.CustomerPointDto) (*entity.CustomerPointDto, error) {
	m, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	m, err := s.repo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) GetTotalPoint(ctx context.Context, customerID string) (float64, error) {
	m, err := s.repo.GetTotalPoint(ctx, customerID)
	if err != nil {
		return 0, err
	}
	return m, nil
}
