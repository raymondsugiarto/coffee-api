package fee_setting

import (
	"context"
	"errors"
	"fmt"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"gorm.io/gorm"
)

const (
	MinFeePercent = 0.0
	MaxFeePercent = 100.0
)

type Service interface {
	UpsertConfig(ctx context.Context, dto *entity.FeeSettingDto) error
	GetConfig(ctx context.Context) (*entity.FeeSettingDto, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) UpsertConfig(ctx context.Context, dto *entity.FeeSettingDto) error {
	if dto.AdminFee < MinFeePercent || dto.AdminFee > MaxFeePercent {
		return status.New(status.InvalidFieldFormat, fmt.Errorf("admin fee must be between %.0f and %.0f percent, got %.2f", MinFeePercent, MaxFeePercent, dto.AdminFee))
	}

	if dto.OperationalFee < MinFeePercent || dto.OperationalFee > MaxFeePercent {
		return status.New(status.InvalidFieldFormat, fmt.Errorf("operational fee must be between %.0f and %.0f percent, got %.2f", MinFeePercent, MaxFeePercent, dto.OperationalFee))
	}

	return s.repo.UpsertConfig(ctx, dto)
}

func (s *service) GetConfig(ctx context.Context) (*entity.FeeSettingDto, error) {
	res, err := s.repo.GetConfig(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.New(status.EntityNotFound, fmt.Errorf("failed to get fee setting: %w", err))
		}
		return nil, status.New(status.InternalServerError, fmt.Errorf("failed to get fee setting: %w", err))
	}
	return res, nil
}
