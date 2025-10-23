package pensionbenefitrecipient

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, recipient *entity.PensionBenefitRecipientDto) (*entity.PensionBenefitRecipientDto, error)
	FindAll(ctx context.Context, req *entity.RecipientFindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.PensionBenefitRecipientDto, error)
	Update(ctx context.Context, recipient *entity.PensionBenefitRecipientDto) (*entity.PensionBenefitRecipientDto, error)
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, customerID string, recipients []*entity.PensionBenefitRecipientDto) error
	CreateWithTx(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error)
	UpdateWithTx(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(ctx context.Context, recipient *entity.PensionBenefitRecipientDto) (*entity.PensionBenefitRecipientDto, error) {
	return s.repository.Create(ctx, recipient, nil)
}

func (s *service) FindAll(ctx context.Context, req *entity.RecipientFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.PensionBenefitRecipientDto, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *service) Update(ctx context.Context, recipient *entity.PensionBenefitRecipientDto) (*entity.PensionBenefitRecipientDto, error) {
	return s.repository.Update(ctx, recipient, nil)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) BatchDelete(ctx context.Context, customerID string, recipients []*entity.PensionBenefitRecipientDto) error {
	return s.repository.BatchDelete(ctx, customerID, recipients)
}

func (s *service) CreateWithTx(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error) {
	return s.repository.Create(ctx, recipient, tx)
}

func (s *service) UpdateWithTx(ctx context.Context, recipient *entity.PensionBenefitRecipientDto, tx *gorm.DB) (*entity.PensionBenefitRecipientDto, error) {
	return s.repository.Update(ctx, recipient, tx)
}
