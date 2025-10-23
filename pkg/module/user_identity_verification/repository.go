package useridentityverification

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error)
	Update(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error)
	FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	m := dto.ToModel()
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	dto.ID = m.ID

	return new(entity.UserIdentityVerificationDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.UserIdentityVerificationDto) (*entity.UserIdentityVerificationDto, error) {
	m := dto.ToModel()
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.UserIdentityVerificationDto, error) {
	m := new(model.UserIdentityVerification)
	if err := r.db.Where("id = ?", id).First(m).Error; err != nil {
		return nil, err
	}

	return new(entity.UserIdentityVerificationDto).FromModel(m), nil
}

func (r *repository) FindByIDAndUniqueCode(ctx context.Context, id, uniqueCode string) (*entity.UserIdentityVerificationDto, error) {
	m := new(model.UserIdentityVerification)
	if err := r.db.Where("id = ? AND unique_code = ?", id, uniqueCode).First(m).Error; err != nil {
		return nil, err
	}
	return new(entity.UserIdentityVerificationDto).FromModel(m), nil
}
