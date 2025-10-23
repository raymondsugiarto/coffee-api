package benefitparticipation

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/model"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"gorm.io/gorm"
)

type CallbackCreate func(ctx context.Context, db *gorm.DB, dto *entity.BenefitParticipationDto) error

type Repository interface {
	Create(ctx context.Context, benefitParticipation *entity.BenefitParticipationDto, cb CallbackCreate) (*entity.BenefitParticipationDto, error)
	CreateWithTx(ctx context.Context, tx *gorm.DB, benefitParticipation *entity.BenefitParticipationDto) (*entity.BenefitParticipationDto, error)
	FindByID(ctx context.Context, id string) (*entity.BenefitParticipationDto, error)
	FindByParticipantID(ctx context.Context, participantID string) (*entity.BenefitParticipationDto, error)
	AppendOrUpdateDetails(ctx context.Context, benefitParticipationID string, newDetails, updateDetails []*entity.BenefitParticipationDetailDto, cb CallbackCreate) error
	UpdateStatus(ctx context.Context, id string, status model.BenefitParticipationStatus, tx *gorm.DB) error
	UpdateWithTx(ctx context.Context, benefitParticipation *entity.BenefitParticipationDto, tx *gorm.DB) error
	FindByInvestmentPaymentID(ctx context.Context, investmentPaymentID string) (*entity.BenefitParticipationDto, error)
	UpdateDetailsStatusByBenefitParticipationID(ctx context.Context, benefitParticipationID string, status model.BenefitParticipationStatus, tx *gorm.DB) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, benefitRegistration *entity.BenefitParticipationDto, cb CallbackCreate) (*entity.BenefitParticipationDto, error) {
	m := benefitRegistration.ToModel()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(m).Error; err != nil {
			return err
		}
		benefitRegistration.ID = m.ID
		// Update details with the parent ID
		for i, detail := range benefitRegistration.Details {
			if detail != nil && i < len(m.Details) && m.Details[i] != nil {
				benefitRegistration.Details[i].ID = m.Details[i].ID
				benefitRegistration.Details[i].BenefitParticipationID = m.ID
			}
		}
		if cb != nil {
			if err := cb(ctx, tx, benefitRegistration); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return benefitRegistration, nil
}

func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, benefitRegistration *entity.BenefitParticipationDto) (*entity.BenefitParticipationDto, error) {
	m := benefitRegistration.ToModel()
	if err := tx.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return benefitRegistration.FromModel(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.BenefitParticipationDto, error) {
	var m model.BenefitParticipation
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Participant").
		Preload("Details").
		First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return new(entity.BenefitParticipationDto).FromModel(&m), nil
}

func (r *repository) FindByParticipantID(ctx context.Context, participantID string) (*entity.BenefitParticipationDto, error) {
	var m model.BenefitParticipation
	if err := r.db.WithContext(ctx).Model(&m).
		Preload("Customer").
		Preload("Participant").
		Preload("Details").
		First(&m, "participant_id = ?", participantID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return new(entity.BenefitParticipationDto).FromModel(&m), nil
}

func (r *repository) AppendOrUpdateDetails(ctx context.Context, benefitParticipationID string, newDetails, updateDetails []*entity.BenefitParticipationDetailDto, cb CallbackCreate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create new details
		for _, detail := range newDetails {
			detailModel := detail.ToModel()
			detailModel.BenefitParticipationID = benefitParticipationID
			if err := tx.WithContext(ctx).Create(detailModel).Error; err != nil {
				return err
			}
			detail.ID = detailModel.ID
			detail.BenefitParticipationID = benefitParticipationID
		}

		// Update existing details
		for _, detail := range updateDetails {
			detailModel := detail.ToModel()
			if err := tx.WithContext(ctx).Model(&model.BenefitParticipationDetail{}).
				Where("id = ?", detail.ID).
				Updates(detailModel).Error; err != nil {
				return err
			}
		}

		if cb != nil {
			benefitParticipation, err := r.FindByID(ctx, benefitParticipationID)
			if err != nil {
				return err
			}
			if err := cb(ctx, tx, benefitParticipation); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status model.BenefitParticipationStatus, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Model(&model.BenefitParticipation{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *repository) UpdateWithTx(ctx context.Context, benefitParticipation *entity.BenefitParticipationDto, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	m := benefitParticipation.ToModel()
	return db.WithContext(ctx).Model(&model.BenefitParticipation{}).
		Where("id = ?", benefitParticipation.ID).
		Updates(m).Error
}

func (r *repository) FindBenefitParticipationDetailByID(ctx context.Context, id string) (*entity.BenefitParticipationDetailDto, error) {
	var m model.BenefitParticipationDetail
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return new(entity.BenefitParticipationDetailDto).FromModel(&m), nil
}

func (r *repository) UpdateBenefitParticipationDetail(ctx context.Context, detail *entity.BenefitParticipationDetailDto, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	m := detail.ToModel()
	return db.WithContext(ctx).Model(&model.BenefitParticipationDetail{}).Where("id = ?", detail.ID).Updates(m).Error
}

func (r *repository) FindByInvestmentPaymentID(ctx context.Context, investmentPaymentID string) (*entity.BenefitParticipationDto, error) {
	var m model.BenefitParticipation
	if err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Participant").
		Preload("Details").
		Joins("JOIN investment ON benefit_participation.investment_id = investment.id").
		Joins("JOIN investment_payment ON investment.id = investment_payment.investment_id").
		Where("investment_payment.id = ?", investmentPaymentID).
		First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return new(entity.BenefitParticipationDto).FromModel(&m), nil
}

func (r *repository) UpdateDetailsStatusByBenefitParticipationID(ctx context.Context, benefitParticipationID string, status model.BenefitParticipationStatus, tx *gorm.DB) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).Model(&model.BenefitParticipationDetail{}).
		Where("benefit_participation_id = ?", benefitParticipationID).
		Update("status", status).Error
}
