package user

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
)

type Service interface {
	// CreateUser(ctx context.Context, req *entity.UserDto) (*entity.UserDto, error)
	FindByReferralCode(ctx context.Context, referralCode string) (*entity.UserDto, error)
	FindByID(ctx context.Context, id string) (*entity.UserDto, error)
	UpdatePhoneVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error
	UpdateEmailVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error
}

type service struct {
	repository            Repository
	userCredentialService usercredential.Service
}

func NewService(repository Repository, userCredentialService usercredential.Service) Service {
	return &service{
		repository:            repository,
		userCredentialService: userCredentialService,
	}
}

func (s *service) UpdatePhoneVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error {
	return s.repository.UpdatePhoneVerificationStatus(ctx, id, status)
}

func (s *service) UpdateEmailVerificationStatus(ctx context.Context, id string, status model.IdentityStatus) error {
	return s.repository.UpdateEmailVerificationStatus(ctx, id, status)
}

func (s *service) FindByReferralCode(ctx context.Context, referralCode string) (*entity.UserDto, error) {
	return s.repository.FindByReferralCode(ctx, referralCode)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.UserDto, error) {
	return s.repository.FindByID(ctx, id)
}

// func (s *service) CreateUser(ctx context.Context, userDto *entity.UserDto) (*entity.UserDto, error) {
// 	userCredential := &entity.UserCredentialDto{
// 		Organization: userDto.OrganizationData,
// 		Username:     userDto.Username,
// 	}
// 	_, err := s.userCredentialService.FindByUsername(ctx, userCredential)
// 	if err == nil {
// 		return nil, errors.New("errorAccountCodeAlreadyExist")
// 	}

// 	_, err = s.userCredentialService.FindByEmail(ctx, userCredential)
// 	if err == nil {
// 		return nil, errors.New("errorEmailAlreadyExist")
// 	}

// 	return s.repository.CreateUser(ctx, userDto)
// }
