package investmentdistribution

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
	userlog "github.com/raymondsugiarto/coffee-api/pkg/module/user_log"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
)

type Service interface {
	CreateBatch(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error)
	Create(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error)
	FindByID(ctx context.Context, id string) (*entity.InvestmentDistributionDto, error)
	Update(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionDto, error)
	FindByParticipantID(ctx context.Context, participantID string) ([]*entity.InvestmentDistributionDto, error)
	CreateOrUpdate(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error)

	HaveDistribution(ctx context.Context, companyID string) error

	SummaryByCompany(ctx context.Context) ([]*entity.InvestmentDistributionSummaryCompanyDto, error)
}

type service struct {
	repo                     Repository
	participantService       participant.Service
	investmentProductService investmentproduct.Service
	userLogService           userlog.Service
}

func NewService(repo Repository, participantService participant.Service, investmentProductService investmentproduct.Service, userlogService userlog.Service) Service {
	return &service{
		repo:                     repo,
		participantService:       participantService,
		userLogService:           userlogService,
		investmentProductService: investmentProductService,
	}
}
func (s *service) CreateBatch(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	companyID := shared.GetCompanyID(ctx)
	for _, item := range dto {
		item.OrganizationID = organizationID
		if companyID != nil {
			item.CompanyID = *companyID
		}
	}

	dtos, err := s.repo.CreateBatch(ctx, dto)
	if err != nil {
		return nil, err
	}

	for _, item := range dtos {
		s.writeUserLog(item, ctx)
	}

	return dtos, nil
}

func (s *service) HaveDistribution(ctx context.Context, companyID string) error {
	dto, err := s.FindByCompanyID(ctx, companyID)
	if err != nil {
		return err
	}
	if len(dto) == 0 {
		return fmt.Errorf("tidak ada alokasi investasi untuk perusahaan %s", companyID)
	}

	return nil
}

func (s *service) writeUserLog(item *entity.InvestmentDistributionDto, ctx context.Context) {
	userLogDto := item.ToUserLogDto()
	investmentProductService, err := s.investmentProductService.FindByID(ctx, item.InvestmentProductID, false)
	log.WithContext(ctx).Infof("Investment Product Service: %v", investmentProductService)
	if err == nil {
		userLogDto.UserCredentialID = shared.GetUserCredential(ctx).ID
		userLogDto.Description = fmt.Sprintf("Alokasi Investasi %s = %.2f%%", investmentProductService.Name, item.Percent)
		s.userLogService.Create(ctx, userLogDto)
	}
}

func (s *service) Create(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	res, err := s.repo.Create(ctx, dto)
	if err != nil {
		return nil, err
	}
	s.writeUserLog(res, ctx)
	return res, nil
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.InvestmentDistributionDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	res, err := s.repo.Update(ctx, dto)
	if err != nil {
		return nil, err
	}
	s.writeUserLog(res, ctx)
	return res, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionDto, error) {
	dto, err := s.repo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	if dto == nil {
		return nil, nil
	}
	return dto, nil
}

func (s *service) FindByParticipantID(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionDto, error) {
	dto, err := s.repo.FindByParticipantID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	if dto == nil {
		return nil, nil
	}
	return dto, nil
}

func (s *service) CreateOrUpdate(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error) {
	organizationID := shared.GetOrganization(ctx).ID
	var (
		dtoCreate []*entity.InvestmentDistributionDto
		result    []*entity.InvestmentDistributionDto
	)

	for _, item := range dto {
		item.OrganizationID = organizationID
		if item.ID == "" {
			dtoCreate = append(dtoCreate, item)
		} else {
			update, err := s.Update(ctx, item)
			if err != nil {
				return nil, err
			}
			result = append(result, update)
		}
	}

	if len(dtoCreate) > 0 {
		createBatch, err := s.CreateBatch(ctx, dtoCreate)
		if err != nil {
			return nil, err
		}
		result = append(result, createBatch...)
	}
	return result, nil
}

func (s *service) SummaryByCompany(ctx context.Context) ([]*entity.InvestmentDistributionSummaryCompanyDto, error) {
	companyID := shared.GetCompanyID(ctx)
	m, err := s.repo.SummaryByCompany(ctx, *companyID)
	if err != nil {
		return nil, err
	}
	return m, nil
}
