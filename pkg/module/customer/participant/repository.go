package participant

import (
	"context"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error)
	Get(ctx context.Context, id string) (*entity.ParticipantDto, error)
	Update(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.ParticipantFindAllRequest) (*pagination.ResultPagination, error)

	FindAllParticipantCompany(ctx context.Context, req *entity.ParticipantFindAllRequest) ([]*entity.ParticipantDto, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ParticipantDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ParticipantDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.ParticipantDto, error) {
	var m *model.Participant

	err := r.db.Where("id = ?", id).Preload("Customer").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.ParticipantDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.ParticipantDto) (*entity.ParticipantDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.Participant{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.ParticipantFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.Participant = make([]model.Participant, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.Participant{}).
			Preload("Customer")

		if req.CompanyID != nil {
			q = q.Joins("JOIN customer ON customer.id = participant.customer_id").
				Joins("JOIN company ON company.id = customer.company_id").
				Preload("Customer.Company").
				Where("company.id = ?", req.CompanyID)
		}

		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"customer_id", "participant.status"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.Participant)
	var data []*entity.ParticipantDto = make([]*entity.ParticipantDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.ParticipantDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindAllParticipantCompany(ctx context.Context, req *entity.ParticipantFindAllRequest) ([]*entity.ParticipantDto, error) {
	var m []model.Participant = make([]model.Participant, 0)

	subquery := r.db.Model(&model.InvestmentItem{}).
		Joins("JOIN investment ON investment.id = investment_item.investment_id").
		Joins("JOIN investment_payment ON investment_payment.investment_id = investment.id").
		Where("investment_payment.status IN ('success','pending')")
	if req.UsePeriod {
		subquery.Where(
			"EXTRACT(YEAR FROM investment.investment_at) = ? AND EXTRACT(MONTH FROM investment.investment_at) = ?",
			req.InvestmentAt.Year(),
			int(req.InvestmentAt.Month()),
		)
	} else {
		subquery.Where("date(investment.investment_at) = ? AND investment.type = 'DKP' AND investment.deleted_at is null", req.InvestmentAt.Format("2006-01-02"))
	}
	subquery.Select("investment_item.participant_id")

	query := r.db.WithContext(ctx).Model(&model.Participant{}).
		Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ? AND customer.approval_status = 'APPROVED' AND customer.sim_status = 'ACTIVE' AND customer.deleted_at is null", req.CompanyID)

	if !req.CalculateAll {
		if req.PaidEmployee {
			query = query.Where("participant.id IN (?)", subquery)
		} else {
			query = query.Where("participant.id NOT IN (?)", subquery)
		}
	}
	err := query.Preload("Customer").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var data []*entity.ParticipantDto = make([]*entity.ParticipantDto, 0)
	for _, m := range m {
		data = append(data, new(entity.ParticipantDto).FromModel(&m))
	}
	return data, nil
}

func (r *repository) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.ParticipantDto, error) {
	var m []model.Participant = make([]model.Participant, 0)

	err := r.db.WithContext(ctx).Joins("JOIN customer ON customer.id = participant.customer_id").
		Where("customer.company_id = ?", companyID).
		Preload("Customer").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	var data []*entity.ParticipantDto = make([]*entity.ParticipantDto, 0)
	for _, m := range m {
		data = append(data, new(entity.ParticipantDto).FromModel(&m))
	}
	return data, nil
}
