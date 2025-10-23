package redeem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	customerpoint "github.com/raymondsugiarto/coffee-api/pkg/module/customer/customer_point"
	"github.com/raymondsugiarto/coffee-api/pkg/module/reward"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"gorm.io/gorm"
)

type Service interface {
	RedeemReward(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error)
	GetAllRedemptions(ctx context.Context, req *entity.RedeemFindAllRequest) (*pagination.ResultPagination, error)
	UpdateRedemptionStatus(ctx context.Context, id string, statusDto *entity.UpdateRedeemStatusDto) (*entity.RedeemDto, error)
	GetRedemptionByID(ctx context.Context, id string) (*entity.RedeemDto, error)
}

type service struct {
	repo              Repository
	rewardRepo        reward.Repository
	customerPointRepo customerpoint.Repository
	customerService   customer.Service
	customerPointSvc  customerpoint.Service
	rewardService     reward.Service
	db                *gorm.DB
}

func NewService(repo Repository, rewardRepo reward.Repository, customerPointRepo customerpoint.Repository, customerService customer.Service, customerPointSvc customerpoint.Service, rewardService reward.Service, db *gorm.DB) Service {
	return &service{
		repo:              repo,
		rewardRepo:        rewardRepo,
		customerPointRepo: customerPointRepo,
		customerService:   customerService,
		customerPointSvc:  customerPointSvc,
		rewardService:     rewardService,
		db:                db,
	}
}

func (s *service) RedeemReward(ctx context.Context, dto *entity.RedeemDto) (*entity.RedeemDto, error) {

	var result *entity.RedeemDto

	err := s.db.Transaction(func(tx *gorm.DB) error {
		rewardDto, err := s.rewardRepo.GetWithLock(ctx, tx, dto.RewardID)
		if err != nil {
			return status.New(status.EntityNotFound, errors.New("reward not found"))
		}

		if rewardDto.Stock <= 0 {
			return status.New(status.BadRequest, errors.New("reward out of stock"))
		}

		totalPoints, err := s.customerPointRepo.GetTotalPointWithLock(ctx, tx, dto.CustomerID)
		if err != nil {
			return err
		}

		if totalPoints < float64(rewardDto.Points) {
			return status.New(status.BadRequest, errors.New("insufficient points"))
		}

		redemptionCode := s.generateRedemptionCode()

		dto.RedemptionDate = time.Now()
		dto.PointsRedeemed = rewardDto.Points
		dto.RedemptionCode = redemptionCode

		createdRedemption, err := s.repo.CreateWithTx(ctx, tx, dto)
		if err != nil {
			return err
		}

		// Deduct points
		pointDeductDto := &entity.CustomerPointDto{
			OrganizationID: dto.OrganizationID,
			CustomerID:     dto.CustomerID,
			Point:          -float64(rewardDto.Points),
			Direction:      "DEBIT",
			Description:    fmt.Sprintf("Redeem reward: %s", rewardDto.Name),
			RefID:          createdRedemption.ID,
			RefModule:      "REDEEM",
			RefCode:        redemptionCode,
		}

		_, err = s.customerPointRepo.CreateWithTx(ctx, tx, pointDeductDto)
		if err != nil {
			return err
		}

		// Decrement stock
		err = s.rewardRepo.DecrementStock(ctx, tx, dto.RewardID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return status.New(status.BadRequest, errors.New("reward out of stock"))
			}
			return err
		}

		result = createdRedemption
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetAllRedemptions(ctx context.Context, req *entity.RedeemFindAllRequest) (*pagination.ResultPagination, error) {
	return s.repo.FindAll(ctx, req)
}

func (s *service) UpdateRedemptionStatus(ctx context.Context, id string, statusDto *entity.UpdateRedeemStatusDto) (*entity.RedeemDto, error) {
	redemption, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.New(status.EntityNotFound, errors.New("redemption not found"))
	}

	oldStatus := redemption.Status
	redemption.Status = statusDto.Status

	var result *entity.RedeemDto

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// If status changed to REJECTED, refund points and restore stock
		if oldStatus != model.RedeemRejected && statusDto.Status == model.RedeemRejected {
			// Refund points
			pointRefundDto := &entity.CustomerPointDto{
				OrganizationID: redemption.OrganizationID,
				CustomerID:     redemption.CustomerID,
				Point:          float64(redemption.PointsRedeemed),
				Direction:      "CREDIT",
				Description:    fmt.Sprintf("Refund for rejected redemption: %s", redemption.RedemptionCode),
				RefID:          redemption.ID,
				RefModule:      "REDEEM",
				RefCode:        redemption.RedemptionCode,
			}

			_, err = s.customerPointRepo.CreateWithTx(ctx, tx, pointRefundDto)
			if err != nil {
				return err
			}

			err = s.rewardRepo.IncrementStock(ctx, tx, redemption.RewardID)
			if err != nil {
				return err
			}
		}

		updatedRedemption, err := s.repo.UpdateWithTx(ctx, tx, redemption)
		if err != nil {
			return err
		}

		result = updatedRedemption
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetRedemptionByID(ctx context.Context, id string) (*entity.RedeemDto, error) {
	redemption, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.New(status.EntityNotFound, errors.New("redemption not found"))
	}
	return redemption, nil
}

func (s *service) generateRedemptionCode() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	randomStr := uuid.New().String()[:8]
	return fmt.Sprintf("RDM-%s-%s", dateStr, randomStr)
}
