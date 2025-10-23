package investmentpayment

import (
	"context"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
	Get(ctx context.Context, id string) (*entity.InvestmentPaymentDto, error)
	Update(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error)
	UpdatePaymentConfirmation(ctx context.Context, dto *entity.InvestmentPaymentDto, cb func(tx *gorm.DB) error) (*entity.InvestmentPaymentDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest) (*pagination.ResultPagination, error)
	GetPaymentSummary(ctx context.Context, companyID *string) (*entity.InvestmentPaymentSummaryDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	m := dto.ToModel()
	err := r.db.Create(m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentPaymentDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.InvestmentPaymentDto, error) {
	var m *model.InvestmentPayment
	err := r.db.Where("id = ?", id).Preload("Investment").Preload("Investment.Company").Preload("Investment.Participant.Customer").First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentPaymentDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.InvestmentPaymentDto) (*entity.InvestmentPaymentDto, error) {
	err := r.db.Updates(dto.ToModel()).Where("id = ? ", dto.ID).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdatePaymentConfirmation(ctx context.Context, dto *entity.InvestmentPaymentDto, cb func(tx *gorm.DB) error) (*entity.InvestmentPaymentDto, error) {
	m := dto.PaymentConfirmationToModel()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.InvestmentPayment{}).
			Where("id = ?", dto.ID).
			Updates(m).Error; err != nil {
			return err
		}

		if err := cb(tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return new(entity.InvestmentPaymentDto).FromModel(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.InvestmentPayment{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.InvestmentPaymentFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.InvestmentPayment = make([]model.InvestmentPayment, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.InvestmentPayment{}).
			Joins("JOIN investment ON investment.id = investment_payment.investment_id").
			Joins("JOIN company ON company.id = investment.company_id").
			Preload("Investment").
			Preload("Investment.Company").
			Preload("Investment.Participant.Customer")

		if req.CompanyID != "" {
			q = q.Where("company.id = ?", req.CompanyID)
		}

		if req.CustomerID != "" {
			q = q.Joins("JOIN customer ON customer.id = investment.customer_id").
				Preload("Investment.Customer").
				Where("customer.id = ?", req.CustomerID)
		}

		return q
	}, &pagination.TableRequest{
		Request: req,
		QueryField: []string{
			"investment.code",
			"company.first_name",
		},
		MapFields: map[string]string{
			"code": "investment.code",
		},
		Data:          &m,
		AllowedFields: []string{"payment_at", "created_at"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.InvestmentPayment)
	var data []*entity.InvestmentPaymentDto = make([]*entity.InvestmentPaymentDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.InvestmentPaymentDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) GetPaymentSummary(ctx context.Context, companyID *string) (*entity.InvestmentPaymentSummaryDto, error) {
	var result entity.InvestmentPaymentSummaryDto

	query := r.db.Model(&model.InvestmentPayment{}).
		Select(`
			COALESCE(SUM(CASE WHEN investment_payment.status IN ('success', 'confirmed') THEN investment_payment.amount ELSE 0 END), 0) as total_payments,
			COALESCE(SUM(CASE WHEN investment_payment.status = 'success' THEN investment_payment.amount ELSE 0 END), 0) as total_approved_payments,
			COALESCE(SUM(CASE WHEN investment_payment.status = 'confirmed' THEN investment_payment.amount ELSE 0 END), 0) as total_pending_payments
		`)

	if companyID != nil {
		query = query.Joins("JOIN investment ON investment.id = investment_payment.investment_id").
			Where("investment.company_id = ?", *companyID)
	}

	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}
