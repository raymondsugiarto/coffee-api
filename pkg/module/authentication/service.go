package authentication

import (
	"context"
	"errors"
	"fmt"

	e "github.com/raymondsugiarto/coffee-api/pkg/entity"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication/token"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
)

type Service interface {
	SignIn(context.Context, *entity.LoginRequestDto) (*entity.LoginDto, error)
}

type service struct {
	userCredentialService usercredential.Service
	tokenService          token.Service
	adminService          admin.Service
}

func NewService(
	userCredentialService usercredential.Service,
	tokenService token.Service,
	adminService admin.Service,
) Service {
	return &service{
		userCredentialService: userCredentialService,
		tokenService:          tokenService,
		adminService:          adminService,
	}
}

func (s *service) SignIn(ctx context.Context, request *entity.LoginRequestDto) (*entity.LoginDto, error) {
	userCredentialDto, err := s.userCredentialService.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	userCredentialData := e.UserCredentialData{
		ID:     userCredentialDto.ID,
		UserID: userCredentialDto.User.ID,
	}

	if userCredentialDto.User.UserType == "ADMIN" {
		admin, err := s.adminService.FindByUserID(ctx, userCredentialDto.User.ID)
		if err != nil {
			return nil, err
		}
		userCredentialData.AdminID = admin.ID
	}

	pp, _ := utils.HashPassword(request.Password)
	fmt.Printf("password hash: %+v\n", pp)
	// fmt.Printf("userCredentialData: %+v ::: %+v\n", request.Password, userCredentialDto.Password)
	if !utils.CheckPasswordHash(request.Password, userCredentialDto.Password) {
		return nil, errors.New("invalidPassword")
	}
	return s.tokenService.GenerateToken(ctx, userCredentialData)
}
