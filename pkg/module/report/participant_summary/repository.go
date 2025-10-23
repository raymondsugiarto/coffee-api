package participantsummary

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetSummaryAll(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error)
	GetSummaryPilar(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error)
	GetSummaryNonPilar(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error)
	GetSummaryManfaatLain(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error)
	GetSummaryPerusahaan(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetSummaryAll(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error) {
	var result []entity.ReportParticipantSummary

	err := r.db.WithContext(ctx).
		Table("customer").
		Select(`
			CASE
				WHEN NULLIF(TRIM(customer.company_id), '') IS NULL THEN 'PPIP'
				WHEN company.company_type = 'PPIP' THEN 'PPIP'
				WHEN company.company_type = 'DKP' THEN 'DKP'
			END as participant_type,
			CASE
				WHEN NULLIF(TRIM(customer.company_id), '') IS NULL THEN 'MANDIRI'
				ELSE 'KUMPULAN'
			END as participant_category,
			'ALL' as pilar_type,
			COUNT(customer.id) as participant_count,
			'ALL' as data_source
		`).
		Joins("LEFT JOIN company ON company.id = customer.company_id").
		Where("customer.deleted_at IS NULL").
		Where("customer.sim_status = ?", model.SIMStatusActive).
		Where("customer.effective_date <= ?", filter.EndDate).
		Group("participant_type, participant_category").
		Scan(&result).Error

	return result, err
}

func (r *repository) GetSummaryPilar(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error) {
	var result []entity.ReportParticipantSummary

	err := r.db.WithContext(ctx).
		Table("customer").
		Select(`
			CASE
				WHEN company.company_type = 'PPIP' THEN 'PPIP'
				WHEN company.company_type = 'DKP' THEN 'DKP'
			END as participant_type,
			'KUMPULAN' as participant_category,
			'PILAR' as pilar_type,
			COUNT(customer.id) as participant_count,
			'PILAR' as data_source
		`).
		Joins("JOIN company ON company.id = customer.company_id").
		Where("customer.deleted_at IS NULL").
		Where("customer.sim_status = ?", model.SIMStatusActive).
		Where("customer.effective_date <= ?", filter.EndDate).
		Where("company.pilar_type = ?", model.TypePilar).
		Group("participant_type").
		Scan(&result).Error

	return result, err
}

func (r *repository) GetSummaryNonPilar(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error) {
	var result []entity.ReportParticipantSummary

	err := r.db.WithContext(ctx).
		Table("customer").
		Select(`
			CASE
				WHEN company.company_type = 'PPIP' THEN 'PPIP'
				WHEN company.company_type = 'DKP' THEN 'DKP'
			END as participant_type,
			'KUMPULAN' as participant_category,
			'NON PILAR' as pilar_type,
			COUNT(customer.id) as participant_count,
			'NON_PILAR' as data_source
		`).
		Joins("JOIN company ON company.id = customer.company_id").
		Where("customer.deleted_at IS NULL").
		Where("customer.sim_status = ?", model.SIMStatusActive).
		Where("customer.effective_date <= ?", filter.EndDate).
		Where("company.pilar_type = ?", model.TypeNonPilar).
		Group("participant_type").
		Scan(&result).Error

	return result, err
}

func (r *repository) GetSummaryManfaatLain(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error) {
	var result []entity.ReportParticipantSummary

	err := r.db.WithContext(ctx).
		Table("benefit_participation_detail bpd").
		Select(`
			bt.name as participant_type,
			CASE
				WHEN NULLIF(TRIM(c.company_id), '') IS NULL THEN 'MANDIRI'
				ELSE 'KUMPULAN'
			END as participant_category,
			'' as pilar_type,
			COUNT(p.id) as participant_count,
			'MANFAAT_LAIN' as data_source
		`).
		Joins("JOIN benefit_type bt ON bt.id = bpd.benefit_type_id").
		Joins("JOIN benefit_participation bp ON bp.id = bpd.benefit_participation_id").
		Joins("JOIN participant p ON p.id = bp.participant_id").
		Joins("JOIN customer c ON c.id = p.customer_id").
		Where("c.deleted_at IS NULL").
		Where("c.sim_status = ?", model.SIMStatusActive).
		Where("c.effective_date <= ?", filter.EndDate).
		Group("bt.name, participant_category").
		Scan(&result).Error

	return result, err
}

func (r *repository) GetSummaryPerusahaan(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.ReportParticipantSummary, error) {
	var result []entity.ReportParticipantSummary

	err := r.db.WithContext(context.Background()).
		Table("company").
		Select(`
			company_type as participant_type,
			COALESCE(pilar_type, '') as participant_category,
			'' as pilar_type,
			COUNT(DISTINCT company.id) as participant_count,
			'PERUSAHAAN' as data_source
		`).
		Where("deleted_at IS NULL").
		Where("status = 'APPROVED'").
		Where("created_at <= ?", filter.EndDate).
		Group("company_type, COALESCE(pilar_type, '')").
		Scan(&result).Error

	return result, err
}
