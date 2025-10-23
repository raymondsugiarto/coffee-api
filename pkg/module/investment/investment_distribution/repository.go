package investmentdistribution

import (
	"context"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	CreateBatch(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error)
	Create(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error)
	Get(ctx context.Context, id string) (*entity.InvestmentDistributionDto, error)
	Update(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionDto, error)
	FindByParticipantID(ctx context.Context, participantID string) ([]*entity.InvestmentDistributionDto, error)

	SummaryByCompany(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionSummaryCompanyDto, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateBatch(ctx context.Context, dto []*entity.InvestmentDistributionDto) ([]*entity.InvestmentDistributionDto, error) {
	m := make([]*model.InvestmentDistribution, len(dto))
	for i, d := range dto {
		m[i] = d.ToModel()
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var result []*entity.InvestmentDistributionDto
	for _, item := range m {
		result = append(result, new(entity.InvestmentDistributionDto).FromModel(item))
	}
	return result, nil
}

func (r *repository) Create(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error) {
	m := dto.ToModel()
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentDistributionDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*entity.InvestmentDistributionDto, error) {
	var m *model.InvestmentDistribution
	err := r.db.Where("id = ?", id).
		Preload("Company").
		Preload("InvestmentProduct").
		Preload("Participant").
		Preload("Customer").
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(entity.InvestmentDistributionDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.InvestmentDistributionDto) (*entity.InvestmentDistributionDto, error) {
	err := r.db.Model(&model.InvestmentDistribution{}).Where("id = ? ", dto.ID).
		Update("organization_id", dto.OrganizationID).
		Update("type", dto.Type).
		Update("company_id", dto.CompanyID).
		Update("participant_id", dto.ParticipantID).
		Update("customer_id", dto.CustomerID).
		Update("investment_product_id", dto.InvestmentProductID).
		Update("percent", dto.Percent).
		Update("base_contribution", dto.BaseContribution).
		Update("voluntary_contribution", dto.VoluntaryContribution).
		Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.InvestmentDistribution{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *entity.InvestmentFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.InvestmentDistribution = make([]model.InvestmentDistribution, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.InvestmentDistribution{}).
			Preload("Company").
			Preload("Customer").
			Preload("Participant").
			Preload("InvestmentProduct").
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"customer_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.InvestmentDistribution)
	var data []*entity.InvestmentDistributionDto = make([]*entity.InvestmentDistributionDto, 0)
	for _, m := range *results {
		data = append(data, new(entity.InvestmentDistributionDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindByCompanyID(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionDto, error) {
	var m []model.InvestmentDistribution = make([]model.InvestmentDistribution, 0)
	err := r.db.Where("company_id = ?", companyID).
		Preload("Company").
		Preload("Customer").
		Preload("Participant").
		Preload("InvestmentProduct").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	result := make([]*entity.InvestmentDistributionDto, len(m))
	for i, item := range m {
		result[i] = new(entity.InvestmentDistributionDto).FromModel(&item)
	}
	return result, nil
}

func (r *repository) FindByParticipantID(ctx context.Context, participantID string) ([]*entity.InvestmentDistributionDto, error) {
	var m []model.InvestmentDistribution = make([]model.InvestmentDistribution, 0)
	err := r.db.Where("participant_id = ?", participantID).
		Preload("Company").
		Preload("Customer").
		Preload("Participant").
		Preload("InvestmentProduct").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	result := make([]*entity.InvestmentDistributionDto, len(m))
	for i, item := range m {
		result[i] = new(entity.InvestmentDistributionDto).FromModel(&item)
	}
	return result, nil
}

func (r *repository) SummaryByCompany(ctx context.Context, companyID string) ([]*entity.InvestmentDistributionSummaryCompanyDto, error) {
	var m []*entity.InvestmentDistributionSummaryCompanyDto
	date := time.Now().Add(7 * time.Hour) // Adjust to Jakarta timezone
	today := date.Format("2006-01-02")

	latestNavSubquery := r.db.WithContext(ctx).Model(&model.NetAssetValue{}).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Where("created_date <= ?", today).
		Where("deleted_at IS NULL")

	err := r.db.WithContext(ctx).Model(&model.InvestmentDistribution{}).
		Joins("LEFT JOIN investment_product ON investment_product.id = investment_distribution.investment_product_id").
		Joins("LEFT JOIN unit_link ON unit_link.investment_product_id = investment_distribution.investment_product_id").
		Joins("LEFT JOIN customer ON customer.id = unit_link.customer_id").
		Joins("LEFT JOIN (?) as nav ON nav.investment_product_id = investment_distribution.investment_product_id AND nav.rn = 1", latestNavSubquery).
		Select("investment_distribution.investment_product_id, investment_product.name as investment_product_name, SUM(ip) as sum_ip, SUM(ip * COALESCE(nav.amount, 0)) as invest_amount, (SUM(ip) * COALESCE(nav.amount, 0)) as total_amount, COALESCE(nav.amount, 0) as nab").
		Where("investment_distribution.company_id = ?", companyID).
		Where("customer.company_id = ?", companyID).
		Group("investment_distribution.investment_product_id, investment_product.name, nav.amount").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}
