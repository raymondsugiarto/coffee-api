package usercredential

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(ctx context.Context, username string) (*entity.UserCredentialDto, error)
	FindByEmail(ctx context.Context, req *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	ChangePassword(ctx context.Context, req *entity.ChangePasswordDto) error
	FindAllByUserID(ctx context.Context, userID string) ([]entity.UserCredentialDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// FindByUsername is a function to find user credential by username
func (r *repository) FindByUsername(ctx context.Context, username string) (*entity.UserCredentialDto, error) {
	userType := shared.GetOriginTypeKey(ctx)
	organization := shared.GetOrganization(ctx)
	var userCredentialModel model.UserCredential
	if err := r.db.Joins("User").
		Where(`"user_credential".username = ? AND "user_credential".organization_id = ? AND "User".user_type = ?`,
			username,
			organization.ID,
			userType,
		).
		First(&userCredentialModel).Error; err != nil {
		return nil, err
	}
	return new(entity.UserCredentialDto).FromModel(&userCredentialModel), nil
}

// FindByUsername is a function to find user credential by username
func (r *repository) FindByEmail(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var userCredentialModel model.UserCredential
	if err := r.db.Joins("User").
		Joins("User.Customer").
		Where("user_credential.username = ? AND user_credential.organization_id = ?", &userCredential.Email, organizationID).
		First(&userCredentialModel).Error; err != nil {
		return nil, err
	}
	// userCredential.CustomerID = userCredentialModel.User.Customer.ID
	userCredential.ID = userCredentialModel.ID
	return new(entity.UserCredentialDto).FromModel(&userCredentialModel), nil
}

func (r *repository) ChangePassword(ctx context.Context, changePassword *entity.ChangePasswordDto) error {
	var userCredentialModel model.UserCredential
	if err := r.db.Where("id = ?", changePassword.UserCredentialID).First(&userCredentialModel).Error; err != nil {
		return err
	}
	if err := r.db.Model(&userCredentialModel).Update("password", changePassword.Password).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAllByUserID(ctx context.Context, userID string) ([]entity.UserCredentialDto, error) {
	var userCredentialModel []model.UserCredential
	if err := r.db.Where("user_id = ?", userID).Find(&userCredentialModel).Error; err != nil {
		return nil, err
	}
	dto := make([]entity.UserCredentialDto, len(userCredentialModel))
	for i, v := range userCredentialModel {
		dto[i] = *new(entity.UserCredentialDto).FromModel(&v)
	}
	return dto, nil
}
