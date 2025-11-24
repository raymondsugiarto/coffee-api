package useridentityverification

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
)

type Callback func(ctx context.Context, dto *entity.UserIdentityVerificationDto) error

type Service interface {
	Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error)
	FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error)
}

type service struct {
	repository Repository
}

func NewService(
	repository Repository,
) Service {
	s := &service{
		repository: repository,
	}
	return s
}

func (s *service) Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error) {
	dto.TryCount = 0

	response, err := s.repository.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error) {
	return s.repository.FindByIDAndUniqueCode(ctx, id, uniqueCode)
}
