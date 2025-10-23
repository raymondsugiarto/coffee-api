package netassetvalue

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	ec "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	"github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	transactionfee "github.com/raymondsugiarto/coffee-api/pkg/module/investment/transaction_fee"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error)
	CreateBatch(ctx context.Context, dto *entity.NetAssetValueBatchDto) (*entity.NetAssetValueBatchDto, error)
	FindByID(ctx context.Context, id string) (*entity.NetAssetValueDto, error)
	FindByInvestmentProductAndDate(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error)
	FindByDate(ctx context.Context, date string) ([]*entity.NetAssetValueDto, error)
	Update(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.NetAssetValueFindAllRequest) (*pagination.ResultPagination, error)

	PublishByInvestment(ctx context.Context, dto []*ec.UnitLinkDto) error
	PublishByNetAssetValue(ctx context.Context, dto []*entity.NetAssetValueDto) error
	ExecuteMonthlyFee(ctx context.Context) error
}

type service struct {
	repo                  Repository
	unitLinkService       unitlink.Service
	feeSettingService     fee_setting.Service
	transactionFeeService transactionfee.Service
}

func NewService(
	repo Repository,
	unitLinkService unitlink.Service,
	feeSettingService fee_setting.Service,
	transactionFeeService transactionfee.Service,
) Service {
	return &service{
		repo:                  repo,
		unitLinkService:       unitLinkService,
		feeSettingService:     feeSettingService,
		transactionFeeService: transactionFeeService,
	}
}
func (s *service) CreateBatch(ctx context.Context, dto *entity.NetAssetValueBatchDto) (*entity.NetAssetValueBatchDto, error) {
	newItems := new(entity.NetAssetValueBatchDto)
	newItems.Items = make([]*entity.NetAssetValueDto, 0)

	updatedItems := make([]*entity.NetAssetValueDto, 0)
	for _, item := range dto.Items {
		netAssetValue, err := s.FindByInvestmentProductAndDate(ctx, item.InvestmentProductID, item.CreatedDate.Format("2006-01-02"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newItems.Items = append(newItems.Items, item)
			} else {
				return nil, err
			}
		}
		item.OrganizationID = shared.GetOrganization(ctx).ID
		if netAssetValue != nil {
			item.ID = netAssetValue.ID
			updatedItems = append(updatedItems, item)
		}
	}
	_, err := s.repo.CreateBatchWithCallback(ctx, newItems, func(tx *gorm.DB) error {
		for _, item := range updatedItems {
			if _, err := s.repo.UpdateWithTx(ctx, item, tx); err != nil {
				return err
			}
		}
		return nil
	})

	// trigger all unit link
	go s.PublishByNetAssetValue(ctx, newItems.Items)

	return nil, err
}

func (s *service) FindByInvestmentProductAndDate(ctx context.Context, investmentProductID string, date string) (*entity.NetAssetValueDto, error) {
	return s.repo.FindByInvestmentProductAndDate(ctx, investmentProductID, date)
}

func (s *service) FindByDate(ctx context.Context, date string) ([]*entity.NetAssetValueDto, error) {
	return s.repo.FindByDate(ctx, date)
}

func (s *service) Create(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Create(ctx, dto)
}

func (s *service) FindByID(ctx context.Context, id string) (*entity.NetAssetValueDto, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, dto *entity.NetAssetValueDto) (*entity.NetAssetValueDto, error) {
	dto.OrganizationID = shared.GetOrganization(ctx).ID
	return s.repo.Update(ctx, dto)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) FindAll(ctx context.Context, req *entity.NetAssetValueFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) PublishByNetAssetValue(ctx context.Context, dto []*entity.NetAssetValueDto) error {
	log.WithContext(ctx).Infof("Net Asset Value PublishByNetAssetValue %v", dto)
	unitLinkDtos, err := s.unitLinkService.FindAllByTransactionDate(ctx, dto[0].CreatedDate)
	if err != nil {
		return err
	}

	err = s.executeNabIp(ctx, dto, unitLinkDtos)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) PublishByInvestment(ctx context.Context, dto []*ec.UnitLinkDto) error {
	if len(dto) == 0 {
		log.WithContext(ctx).Info("No unit link dto to publish")
		return nil
	}
	log.WithContext(ctx).Infof("Net Asset Value PublishByInvestment %v", dto)
	netAssetValueDtos, err := s.FindByDate(ctx, dto[0].TransactionDate.Format("2006-01-02"))
	if err != nil {
		return err
	}

	if len(netAssetValueDtos) > 0 {
		return s.executeNabIp(ctx, netAssetValueDtos, dto)
	}
	return nil
}

func (s *service) executeNabIp(ctx context.Context, netAssetValueDtos []*entity.NetAssetValueDto, dto []*ec.UnitLinkDto) error {
	// convert netAssetValueDtos to map[id]netAssetValueDtos using lo.Map
	netAssetValueMap := lo.KeyBy(netAssetValueDtos, func(n *entity.NetAssetValueDto) string {
		return n.InvestmentProductID
	})

	log.WithContext(ctx).Infof("net asset value %v", netAssetValueMap)
	for _, unitLink := range dto {
		unitLink.Nab = netAssetValueMap[unitLink.InvestmentProductID].Amount
		unitLink.Ip = unitLink.TotalAmount / unitLink.Nab
		_, err := s.unitLinkService.UpdateNab(ctx, unitLink)
		if err != nil {
			log.WithContext(ctx).Errorf("error update nab %s", err.Error())
			return err
		}
	}
	return nil
}

func (s *service) ExecuteMonthlyFee(ctx context.Context) error {
	unitLinkPortfolios, err := s.unitLinkService.FindAllInvestmentProductGroupParticipant(ctx)
	if err != nil {
		return err
	}

	unitLinks, err := s.unitLinkService.FindLatestEachProductAndParticipantAndType(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("error find latest each product and participant and type %s", err.Error())
		return err
	}
	unitLinksMap := lo.KeyBy(unitLinks, func(u *ec.UnitLinkLatestEachProductAndParticipantAndTypeDto) string {
		return u.InvestmentProductID + u.ParticipantID + u.Type
	})

	netAssetValues, err := s.FindByDate(ctx, time.Now().Format("2006-01-02"))
	if err != nil {
		return err
	}

	feeSetting, err := s.feeSettingService.GetConfig(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("error get fee setting %s", err.Error())
		return err
	}

	// covert netAssetValues to map using lo.KeyBy
	netAssetValueMap := lo.KeyBy(netAssetValues, func(n *entity.NetAssetValueDto) string {
		return n.InvestmentProductID
	})

	for _, item := range unitLinkPortfolios {
		netAssetValue, ok := netAssetValueMap[item.InvestmentProductID]
		if !ok {
			// should not happen, but just in case
			log.WithContext(ctx).Errorf("net asset value not found for investment product %s", item.InvestmentProductID)
			continue
		}

		unitLink, ok := unitLinksMap[item.InvestmentProductID+item.ParticipantID+item.Type]
		if !ok {
			// should not happen, but just in case
			log.WithContext(ctx).Errorf("unit link not found for investment product %s, participant %s, type %s", item.InvestmentProductID, item.ParticipantID, item.Type)
			continue
		}

		portfolioAmount := item.Ip * netAssetValue.Amount
		operationFeeAmount := portfolioAmount * (feeSetting.OperationalFee / 100)

		// new portfolio amount
		portfolioAmount -= operationFeeAmount
		// new IP
		portfolioIp := portfolioAmount / netAssetValue.Amount
		// diff IP
		diff := item.Ip - portfolioIp

		unitLink.Ip -= diff
		// update latest unit link

		transactionFee := new(entity.TransactionFeeDto)
		transactionFee.InvestmentProductID = item.InvestmentProductID
		transactionFee.ParticipantID = item.ParticipantID
		transactionFee.Type = model.InvestmentType(item.Type)
		transactionFee.TransactionDate = time.Now()
		transactionFee.OrganizationID = item.OrganizationID
		transactionFee.Nav = netAssetValue.Amount
		transactionFee.Ip = unitLink.Ip
		transactionFee.PortfolioAmount = portfolioAmount
		transactionFee.OperationFee = operationFeeAmount

		if _, err := s.transactionFeeService.Create(ctx, transactionFee); err != nil {
			log.WithContext(ctx).Errorf("error create transaction fee %s", err.Error())
			return err
		}

		unitLinkDto := &ec.UnitLinkDto{
			ID:  unitLink.ID,
			Nab: netAssetValue.Amount,
			Ip:  unitLink.Ip,
		}

		_, err := s.unitLinkService.UpdateNab(ctx, unitLinkDto)
		if err != nil {
			return err
		}
	}

	return nil
}
