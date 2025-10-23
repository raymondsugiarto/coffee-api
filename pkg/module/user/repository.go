package user

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"

	"gorm.io/gorm"
)

type Repository interface {
	// Create(ctx context.Context, dto *entity.CreateUser) (*entity.CreateUser, error)
	FindByReferralCode(ctx context.Context, referralCode string) (*entity.UserDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserDto, error)
	UpdatePhoneVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error
	UpdateEmailVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByReferralCode(ctx context.Context, referralCode string) (*entity.UserDto, error) {
	var user *model.User
	if err := r.db.Joins("Customer").
		Where("customer.referral_code = ?", referralCode).
		First(&user).Error; err != nil {
		return nil, err
	}
	return new(entity.UserDto).FromModel(user), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.UserDto, error) {
	var user *model.User
	if err := r.db.
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return new(entity.UserDto).FromModel(user), nil
}

func (r *repository) UpdatePhoneVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error {
	err := r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("phone_verification_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) UpdateEmailVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error {
	err := r.db.Model(&model.User{}).
		Where("id = ?", id).
		Update("email_verification_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateUser(createUser *entity.CreateUser) (*entity.CreateUser, error) {
	// user := new(model.User)
	// user.OrganizationID = createUser.OrganizationData.ID
	// user.Customer = model.Customer{
	// 	OrganizationID: createUser.OrganizationData.ID,
	// 	Email:          createUser.Email,
	// 	FirstName:      createUser.Name,
	// 	PhoneNumber:    createUser.PhoneNumber,
	// }

	// password, _ := utils.HashPassword(createUser.Password)
	// user.UserCredential = []model.UserCredential{
	// 	{
	// 		OrganizationID: createUser.OrganizationData.ID,
	// 		Username:       createUser.Username,
	// 		Password:       password,
	// 	},
	// }

	// if err := r.db.Create(user).Error; err != nil {
	// 	return nil, err
	// }
	// createUser.UserID = user.ID
	// return createUser, nil
	return nil, nil
}
