package admin

import (
	"context"
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, admin *entity.AdminDto) (*entity.AdminDto, error)
	CreateWithTx(ctx context.Context, admin *entity.AdminDto, tx *gorm.DB) (*entity.AdminDto, error)
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	FindByUserID(ctx context.Context, id string) (*entity.AdminDto, error)
	UpdateProfileImage(ctx context.Context, id string, profileImageUrl string) error
	UpdateName(ctx context.Context, id string, admin *entity.AdminDto) (*entity.AdminDto, error)
	CreateAdminCompany(ctx context.Context, admin *entity.CreateAdminCompany) (*entity.CreateAdminCompany, error)
	FindAllByCompanyID(ctx context.Context, companyID string, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
}

type service struct {
	repo               Repository
	roleDefaultCompany string
}

func NewService(repo Repository, roleDefaultCompany string) Service {
	return &service{repo, roleDefaultCompany}
}

func (s *service) UpdateProfileImage(ctx context.Context, id string, profileImageUrl string) error {
	admin, err := s.repo.FindByUserID(ctx, id)
	if err != nil {
		return err
	}

	if admin == nil {
		return status.New(status.EntityNotFound, errors.New("admin not found"))
	}

	admin.ProfileImageUrl = profileImageUrl

	if _, err := s.repo.Update(ctx, admin); err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateName(ctx context.Context, id string, dto *entity.AdminDto) (*entity.AdminDto, error) {
	admin, err := s.repo.FindByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	if admin == nil {
		return nil, status.New(status.EntityNotFound, errors.New("admin not found"))
	}

	admin.FirstName = dto.FirstName
	admin.LastName = dto.LastName

	adminDto, err := s.repo.Update(ctx, admin)
	if err != nil {
		return nil, err
	}

	return adminDto, nil
}

func (s *service) Create(ctx context.Context, admin *entity.AdminDto) (*entity.AdminDto, error) {
	return s.repo.Create(ctx, admin, nil)
}

func (s *service) CreateWithTx(ctx context.Context, admin *entity.AdminDto, tx *gorm.DB) (*entity.AdminDto, error) {
	return s.repo.Create(ctx, admin, tx)
}

func (s *service) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) FindByUserID(ctx context.Context, id string) (*entity.AdminDto, error) {
	return s.repo.FindByUserID(ctx, id)
}

func (s *service) CreateAdminCompany(ctx context.Context, admin *entity.CreateAdminCompany) (*entity.CreateAdminCompany, error) {
	organizationID := shared.GetOrganization(ctx).ID

	hashPassword, _ := utils.HashPassword(admin.Password)
	user := &entity.UserDto{
		OrganizationID: organizationID,
		UserType:       entity.COMPANY,
		UserCredential: []entity.UserCredentialDto{
			{
				OrganizationID: organizationID,
				Username:       admin.Email,
				Password:       hashPassword,
			},
		},
	}

	if admin.OrganizationID == "" {
		admin.OrganizationID = organizationID
	}
	if admin.User == nil {
		admin.User = user
	}
	admin.User.UserHasRoleDto = []entity.UserHasRoleDto{
		{
			RoleID: s.roleDefaultCompany,
		},
	}
	admin.User.EmailVerificationStatus = "UNVERIFIED"
	admin.User.PhoneVerificationStatus = "UNVERIFIED"
	admin.AdminType = string(entity.COMPANY)

	existing, err := s.repo.FindByUserID(ctx, admin.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, status.New(status.BadRequest, errors.New("admin company sudah terdaftar"))
	}

	return s.repo.CreateAdminCompany(ctx, admin, func(tx *gorm.DB) error {

		return nil
	})
}

func (s *service) FindAllByCompanyID(ctx context.Context, companyID string, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAllByCompanyID(ctx, companyID, req)
}
