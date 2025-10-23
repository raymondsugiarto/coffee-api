package company

import (
	"context"

	"github.com/gofiber/fiber/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	ea "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, req *entity.CompanyDto) (*entity.CompanyDto, error)
	FindByEmail(ctx context.Context, email string) (*entity.CompanyDto, error)
	FindByUserID(ctx context.Context, userId string) (*entity.CompanyDto, error)
	FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error)
	FindByID(ctx context.Context, id string) (*entity.CompanyDto, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, dto *entity.CompanyDto) (*entity.CompanyDto, error)
	ConfirmationApprovalCallback(ctx context.Context, req *entity.CompanyDto, tx *gorm.DB) (context.Context, error)
	FindByCompanyCode(ctx context.Context, companyCode string) (*entity.CompanyDto, error)
	CountByType(ctx context.Context) ([]*entity.CountCompanyPerType, error)
	ChangePassword(ctx context.Context, dto *ea.UserIdentityVerificationDto) error
}

type service struct {
	repository            Repository
	adminService          admin.Service
	roleDefaultCompany    string
	approvalService       approval.Service
	userCredentialService usercredential.Service
}

func NewService(repository Repository, adminService admin.Service, roleDefaultCompany string, approvalService approval.Service, userCredentialService usercredential.Service) Service {
	return &service{
		repository:            repository,
		adminService:          adminService,
		roleDefaultCompany:    roleDefaultCompany,
		approvalService:       approvalService,
		userCredentialService: userCredentialService,
	}
}

func (s *service) ConfirmationApprovalCallback(ctx context.Context, req *entity.CompanyDto, tx *gorm.DB) (context.Context, error) {
	// TODO: implement confirmation callback
	// 1. Find company by ID
	company, err := s.repository.FindByID(ctx, req.ID)
	if err != nil {
		return ctx, err
	}
	// 2. Check if company is already confirmed
	if company.Status != req.Status {
		s.repository.Update(ctx, req)
		return ctx, nil
	}
	// 3. If not, update company status to confirmed

	return ctx, nil
}

func (s *service) Create(ctx context.Context, dto *entity.CompanyDto) (*entity.CompanyDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	userType := shared.GetOriginTypeKey(ctx)
	companyCode, err := gonanoid.Generate(utils.ALPHA_NUMERIC, 5)
	if err != nil {
		log.Errorf("errorGenerateCompanyCode: %v", err)
		return nil, status.New(status.BadRequest, err)
	}
	user := &entity.UserDto{
		OrganizationID: organizationID,
		UserType:       entity.UserType(userType),
		UserCredential: []entity.UserCredentialDto{
			{
				OrganizationID: organizationID,
				Username:       dto.Email,
			},
		},
	}

	dto.CompanyCode = companyCode
	if dto.OrganizationID == "" {
		dto.OrganizationID = organizationID
	}
	if dto.User == nil {
		dto.User = user
	}
	dto.User.UserHasRoleDto = []entity.UserHasRoleDto{
		{
			RoleID: s.roleDefaultCompany,
		},
	}

	return s.repository.Create(ctx, dto, func(tx *gorm.DB) error {
		var uid string
		if userCredential := shared.GetUserCredential(ctx); userCredential == nil {
			uid = dto.UserID
		} else {
			uid = userCredential.UserID
		}

		if _, err := s.adminService.CreateWithTx(ctx, dto.ToAdminDto(), tx); err != nil {
			return err
		}
		if _, err := s.approvalService.CreateWithTx(ctx, dto.ToSubmitApprovalDto(uid), tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *service) FindByEmail(ctx context.Context, email string) (*entity.CompanyDto, error) {
	return s.repository.FindByEmail(ctx, email)
}

func (s *service) FindByUserID(ctx context.Context, userId string) (*entity.CompanyDto, error) {
	return s.repository.FindByUserID(ctx, userId)
}

func (s *service) FindAll(ctx context.Context, req *entity.FindAllRequest) (*pagination.ResultPagination, error) {
	return s.repository.FindAll(ctx, req)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.CompanyDto, error) {
	company, err := s.repository.FindByID(ctx, id)
	if err != nil {
		log.Errorf("errorFindByID: %v", err)
		return nil, status.New(status.EntityNotFound, err)
	}

	approval, err := s.approvalService.FindByRefID(ctx, id, "COMPANY")
	if err == nil {
		company.Approval = approval
	}
	return company, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.CompanyDto) (*entity.CompanyDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repository.Update(ctx, dto)
}

func (s *service) FindByCompanyCode(ctx context.Context, companyCode string) (*entity.CompanyDto, error) {
	return s.repository.FindByCompanyCode(ctx, companyCode)
}

func (s *service) CountByType(ctx context.Context) ([]*entity.CountCompanyPerType, error) {
	return s.repository.CountByType(ctx)
}

func (s *service) ChangePassword(ctx context.Context, req *ea.UserIdentityVerificationDto) error {
	userCredentials, err := s.userCredentialService.FindAllByUserID(ctx, req.UserID)
	if err != nil {
		log.Errorf("errorFindAllByUserID: %v", err)
		return err
	}
	data := req.Data.(*ea.UserIdentityVerificationInputPasswordDto)

	for _, userCredential := range userCredentials {
		cp := new(entity.ChangePasswordDto)
		cp.UserCredentialID = userCredential.ID
		cp.Password = data.Password
		if err := s.userCredentialService.ChangePassword(ctx, cp); err != nil {
			return err
		}
	}
	return nil
}
