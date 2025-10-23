package usercredential

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
)

type Service interface {
	FindByUsername(ctx context.Context, username string) (*entity.UserCredentialDto, error)
	FindByEmail(ctx context.Context, req *entity.UserCredentialDto) (*entity.UserCredentialDto, error)
	ChangePassword(ctx context.Context, req *entity.ChangePasswordDto) error
	FindAllByUserID(ctx context.Context, userID string) ([]entity.UserCredentialDto, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

// FindByUsername is a function to find user credential by username
func (s *service) FindByUsername(ctx context.Context, username string) (*entity.UserCredentialDto, error) {
	return s.repository.FindByUsername(ctx, username)
}

// FindByEmail is a function to find user credential by username
func (s *service) FindByEmail(ctx context.Context, userCredential *entity.UserCredentialDto) (*entity.UserCredentialDto, error) {
	return s.repository.FindByEmail(ctx, userCredential)
}

func (s *service) ChangePassword(ctx context.Context, changePassword *entity.ChangePasswordDto) error {
	hashPassword, _ := utils.HashPassword(changePassword.Password)
	changePassword.Password = hashPassword
	return s.repository.ChangePassword(ctx, changePassword)
}

func (s *service) FindAllByUserID(ctx context.Context, userID string) ([]entity.UserCredentialDto, error) {
	return s.repository.FindAllByUserID(ctx, userID)
}
