package unitlink

import (
	"context"
	"errors"
	"time"

	ce "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	CreateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	Get(ctx context.Context, id string) (*ce.UnitLinkDto, error)
	Update(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	UpdateTotalAmountNabIpWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	UpdateNab(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *ce.UnitLinkFindAllRequest) (*pagination.ResultPagination, error)
	FindAllByTransactionDate(ctx context.Context, transactionDate time.Time) ([]*ce.UnitLinkDto, error)
	FindAllInvestmentProductByCustomer(ctx context.Context, customerID string) ([]*ce.UnitLinkPortfolioDto, error)
	FindAllInvestmentProductByParticipant(ctx context.Context, participantID string) ([]*ce.UnitLinkPortfolioDto, error)
	SumInvestmentProductByCustomer(ctx context.Context, customerID string) (*ce.SumUnitLinkPortfolioDto, error)
	SumInvestmentProductByParticipant(ctx context.Context, participantID string) (*ce.SumUnitLinkPortfolioDto, error)
	FindAllInvestmentProductGroupParticipant(ctx context.Context) ([]*ce.UnitLinkPortfolioGroupParticipantDto, error)
	FindLatestEachProductAndParticipantAndType(ctx context.Context) ([]*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto, error)
	SummaryByCompany(ctx context.Context, companyID string) (*ce.UnitLinkSummaryCompanyDto, error)
	SummaryPerType(ctx context.Context) ([]*ce.UnitLinkSummaryPerTypeDto, error)
	FindAllPortfolioWithNav(ctx context.Context, req *ce.PortfolioFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) CreateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m := dto.ToModel()

	existing := &model.UnitLink{}
	err := tx.WithContext(ctx).Where(model.UnitLink{
		TransactionDate:     m.TransactionDate,
		CustomerID:          m.CustomerID,
		ParticipantID:       m.ParticipantID,
		InvestmentProductID: m.InvestmentProductID,
		Type:                m.Type,
	}).First(existing).Error

	isExist := true
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isExist = false
		} else {
			return nil, err
		}
	}

	if isExist {
		m.ID = existing.ID
		m.TotalAmount += existing.TotalAmount
		err = tx.WithContext(ctx).Where("id = ?", existing.ID).Updates(m).Error
	} else {
		err = tx.WithContext(ctx).Create(m).Error
	}

	if err != nil {
		return nil, err
	}
	return new(ce.UnitLinkDto).FromModel(m), nil
}

func (r *repository) Create(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m := dto.ToModel()
	err := r.db.WithContext(ctx).Where(model.UnitLink{
		TransactionDate:     m.TransactionDate,
		ParticipantID:       m.ParticipantID,
		InvestmentProductID: m.InvestmentProductID,
	}).FirstOrCreate(m).Error
	if err != nil {
		return nil, err
	}
	return new(ce.UnitLinkDto).FromModel(m), nil
}

func (r *repository) Get(ctx context.Context, id string) (*ce.UnitLinkDto, error) {
	var m *model.UnitLink
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return new(ce.UnitLinkDto).FromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	err := r.db.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateTotalAmountNabIpWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	m := dto.ToModel()

	updateData := map[string]interface{}{
		"TotalAmount": dto.TotalAmount,
		"Nab":         dto.Nab,
		"Ip":          dto.Ip,
	}

	err := tx.WithContext(ctx).Model(&model.UnitLink{}).Where("id = ?", m.ID).Updates(updateData).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateWithTx(ctx context.Context, tx *gorm.DB, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	err := tx.Updates(dto.ToModel()).Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) UpdateNab(ctx context.Context, dto *ce.UnitLinkDto) (*ce.UnitLinkDto, error) {
	err := r.db.Model(&model.UnitLink{
		CommonWithIDs: concern.CommonWithIDs{
			ID: dto.ID,
		},
	}).Updates(model.UnitLink{
		Nab: dto.Nab,
		Ip:  dto.Ip,
	}).
		Error
	if err != nil {
		return nil, err
	}
	return dto, nil
}
func (r *repository) Delete(ctx context.Context, id string) error {
	err := r.db.Where("id = ?", id).Delete(&model.UnitLink{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context, req *ce.UnitLinkFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.UnitLink = make([]model.UnitLink, 0)

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Model(&model.UnitLink{}).
			Preload("Customer").
			Preload("InvestmentProduct").
			Preload("Participant").
			Where("deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"participant_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.UnitLink)
	var data []*ce.UnitLinkDto = make([]*ce.UnitLinkDto, 0)
	for _, m := range *results {
		data = append(data, new(ce.UnitLinkDto).FromModel(&m))
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}

func (r *repository) FindAllByTransactionDate(ctx context.Context, transactionDate time.Time) ([]*ce.UnitLinkDto, error) {
	var m []model.UnitLink = make([]model.UnitLink, 0)
	err := r.db.WithContext(ctx).Where("transaction_date = ?", transactionDate).Find(&m).Error
	if err != nil {
		return nil, err
	}
	var data []*ce.UnitLinkDto = make([]*ce.UnitLinkDto, 0)
	for _, m := range m {
		data = append(data, new(ce.UnitLinkDto).FromModel(&m))
	}
	return data, nil
}

func (r *repository) FindAllInvestmentProductByCustomer(ctx context.Context, customerID string) ([]*ce.UnitLinkPortfolioDto, error) {
	var m []*ce.UnitLinkPortfolioDto = make([]*ce.UnitLinkPortfolioDto, 0)

	// Calculate IP with fallback to transaction date NAV
	unitWithCalculatedIPSubquery := r.db.WithContext(ctx).
		Select(`ul.investment_product_id,
			CASE
				WHEN ul.ip > 0 THEN ul.ip
				WHEN nav_tx.amount > 0 THEN ul.total_amount / nav_tx.amount
				ELSE 0
			END as calculated_ip`).
		Table("unit_link ul").
		Joins("LEFT JOIN net_asset_value nav_tx ON nav_tx.investment_product_id = ul.investment_product_id AND nav_tx.created_date = ul.transaction_date").
		Where("ul.customer_id = ?", customerID)

	// Get unit summary per product using calculated IP
	unitSummarySubquery := r.db.WithContext(ctx).
		Select("investment_product_id, SUM(calculated_ip) as ip").
		Table("(?) as ucip", unitWithCalculatedIPSubquery).
		Group("investment_product_id")

	// Get latest NAV per product
	latestNavSubquery := r.db.WithContext(ctx).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Table("net_asset_value")

	err := r.db.WithContext(ctx).
		Table("(?) as us", unitSummarySubquery).
		Joins("LEFT JOIN (?) as ln ON ln.investment_product_id = us.investment_product_id AND ln.rn = 1", latestNavSubquery).
		Select("us.investment_product_id, us.ip, (us.ip * COALESCE(ln.amount, 0)) as total_amount, COALESCE(ln.amount, 0) as nab").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) FindAllInvestmentProductByParticipant(ctx context.Context, participantID string) ([]*ce.UnitLinkPortfolioDto, error) {
	var m []*ce.UnitLinkPortfolioDto = make([]*ce.UnitLinkPortfolioDto, 0)

	// Calculate IP with fallback to transaction date NAV
	unitWithCalculatedIPSubquery := r.db.WithContext(ctx).
		Select(`ul.investment_product_id,
			CASE
				WHEN ul.ip > 0 THEN ul.ip
				WHEN nav_tx.amount > 0 THEN ul.total_amount / nav_tx.amount
				ELSE 0
			END as calculated_ip`).
		Table("unit_link ul").
		Joins("LEFT JOIN net_asset_value nav_tx ON nav_tx.investment_product_id = ul.investment_product_id AND nav_tx.created_date = ul.transaction_date").
		Where("ul.participant_id = ?", participantID)

	// Get unit summary per product using calculated IP
	unitSummarySubquery := r.db.WithContext(ctx).
		Select("investment_product_id, SUM(calculated_ip) as ip").
		Table("(?) as ucip", unitWithCalculatedIPSubquery).
		Group("investment_product_id")

	// Get latest NAV per product
	latestNavSubquery := r.db.WithContext(ctx).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Table("net_asset_value")

	err := r.db.WithContext(ctx).
		Table("(?) as us", unitSummarySubquery).
		Joins("LEFT JOIN (?) as ln ON ln.investment_product_id = us.investment_product_id AND ln.rn = 1", latestNavSubquery).
		Select("us.investment_product_id, us.ip, (us.ip * COALESCE(ln.amount, 0)) as total_amount, COALESCE(ln.amount, 0) as nab").
		Find(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) FindAllInvestmentProductGroupParticipant(ctx context.Context) ([]*ce.UnitLinkPortfolioGroupParticipantDto, error) {
	var m []*ce.UnitLinkPortfolioGroupParticipantDto = make([]*ce.UnitLinkPortfolioGroupParticipantDto, 0)

	err := r.db.WithContext(ctx).Model(&model.UnitLink{}).
		Group("unit_link.investment_product_id, unit_link.participant_id, unit_link.type").
		Select("unit_link.investment_product_id, unit_link.participant_id, SUM(ip) as ip").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) FindLatestEachProductAndParticipantAndType(ctx context.Context) ([]*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto, error) {
	var m []*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto = make([]*ce.UnitLinkLatestEachProductAndParticipantAndTypeDto, 0)

	subquery := r.db.WithContext(ctx).Model(&model.UnitLink{}).
		Select("id, investment_product_id, participant_id, type, transaction_date, ip, row_number() over (partition by investment_product_id, participant_id, type order by transaction_date) as rn").
		Where("deleted_at IS NULL")

	err := r.db.WithContext(ctx).Table("(?) as aa", subquery).
		Select("id, investment_product_id, participant_id, type, transaction_date, ip").
		Where("rn = 1").
		Find(&m).Error
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *repository) SumInvestmentProductByCustomer(ctx context.Context, customerID string) (*ce.SumUnitLinkPortfolioDto, error) {
	var m *ce.SumUnitLinkPortfolioDto

	// Calculate IP with fallback to transaction date NAV
	unitWithCalculatedIPSubquery := r.db.WithContext(ctx).
		Select(`ul.investment_product_id,
			CASE
				WHEN ul.ip > 0 THEN ul.ip
				WHEN nav_tx.amount > 0 THEN ul.total_amount / nav_tx.amount
				ELSE 0
			END as calculated_ip,
			ul.total_amount`).
		Table("unit_link ul").
		Joins("LEFT JOIN net_asset_value nav_tx ON nav_tx.investment_product_id = ul.investment_product_id AND nav_tx.created_date = ul.transaction_date").
		Where("ul.customer_id = ?", customerID)

	// Get unit summary per product using calculated IP
	unitSummarySubquery := r.db.WithContext(ctx).
		Select("investment_product_id, SUM(calculated_ip) as total_ip, SUM(total_amount) as modal_per_product").
		Table("(?) as ucip", unitWithCalculatedIPSubquery).
		Group("investment_product_id")

	// Get latest NAV per product
	latestNavSubquery := r.db.WithContext(ctx).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn").
		Table("net_asset_value")

	err := r.db.WithContext(ctx).
		Table("(?) as us", unitSummarySubquery).
		Joins("LEFT JOIN (?) as ln ON ln.investment_product_id = us.investment_product_id AND ln.rn = 1", latestNavSubquery).
		Select("COALESCE(SUM(us.total_ip * COALESCE(ln.amount, 0)), 0) as current_balance, COALESCE(SUM(us.total_ip), 0) as total_unit, COALESCE(SUM(us.modal_per_product), 0) as total_modal, COUNT(*) as count").
		Take(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) SumInvestmentProductByParticipant(ctx context.Context, participantID string) (*ce.SumUnitLinkPortfolioDto, error) {
	var m *ce.SumUnitLinkPortfolioDto

	subquery := r.db.WithContext(ctx).Model(&model.UnitLink{}).
		Preload("InvestmentProduct").
		Joins("LEFT JOIN net_asset_value ON net_asset_value.investment_product_id = unit_link.investment_product_id AND unit_link.transaction_date = net_asset_value.created_date").
		Where("participant_id = ?", participantID).
		Group("unit_link.investment_product_id").
		Select("unit_link.investment_product_id, SUM(ip) as ip, max(created_date) as created_date")

	subquerySummary := r.db.WithContext(ctx).Table("(?) as aa", subquery).
		Joins("LEFT JOIN net_asset_value ON net_asset_value.investment_product_id = aa.investment_product_id AND aa.created_date = net_asset_value.created_date").
		Select("(aa.ip * net_asset_value.amount) as total_amount")

	err := r.db.WithContext(ctx).Table("(?) as aa", subquerySummary).
		Select("sum(total_amount) as total_amount, count(1) as count").
		Take(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) SummaryByCompany(ctx context.Context, companyID string) (*ce.UnitLinkSummaryCompanyDto, error) {
	var m *ce.UnitLinkSummaryCompanyDto

	subquery := r.db.WithContext(ctx).Model(&model.UnitLink{}).
		Joins("JOIN customer ON customer.id = unit_link.customer_id").
		Where("customer.company_id = ?", companyID).
		Group("unit_link.investment_product_id").
		Select("unit_link.investment_product_id, SUM(unit_link.ip) as total_ip")

	latestNavSubquery := r.db.WithContext(ctx).Model(&model.NetAssetValue{}).
		Select("investment_product_id, amount, ROW_NUMBER() OVER (PARTITION BY investment_product_id ORDER BY created_date DESC) as rn")

	subquerySummary := r.db.WithContext(ctx).Table("(?) as ul", subquery).
		Joins("LEFT JOIN (?) as nav ON nav.investment_product_id = ul.investment_product_id AND nav.rn = 1", latestNavSubquery).
		Select("(ul.total_ip * COALESCE(nav.amount, 0)) as total_amount, ul.total_ip as ip")

	err := r.db.WithContext(ctx).Table("(?) as aa", subquerySummary).
		Select("sum(total_amount) as total_amount, count(1) as count, sum(ip) as sum_ip").
		Take(&m).Error
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) SummaryPerType(ctx context.Context) ([]*ce.UnitLinkSummaryPerTypeDto, error) {
	var m []*ce.UnitLinkSummaryPerTypeDto

	date := time.Now().UTC()
	today := date.Format("2006-01-02")

	subquery := r.db.WithContext(ctx).Model(&model.UnitLink{}).
		Preload("InvestmentProduct").
		Joins("JOIN customer ON customer.id = unit_link.customer_id").
		Joins("LEFT JOIN company ON company.id = customer.company_id").
		Group("company.company_type, unit_link.investment_product_id").
		Select("company.company_type, unit_link.investment_product_id, SUM(total_amount) as total_amount, SUM(ip) as ip")

	err := r.db.WithContext(ctx).Table("(?) as aa", subquery).
		Joins("LEFT JOIN net_asset_value ON net_asset_value.investment_product_id = aa.investment_product_id AND net_asset_value.created_date = ?", today).
		Select("company_type as type, SUM(total_amount) as total_amount_unit_link, SUM(aa.ip * net_asset_value.amount) as total_amount, SUM(aa.ip) as sum_ip").
		Group("company_type").
		Find(&m).Error

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) FindAllPortfolioWithNav(ctx context.Context, req *ce.PortfolioFindAllRequest) (*pagination.ResultPagination, error) {
	var m []*ce.PortfolioWithNavDto = make([]*ce.PortfolioWithNavDto, 0)

	date := time.Now().UTC()
	today := date.Format("2006-01-02")

	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		return r.db.Table("unit_link").
			Select(`
				unit_link.id,
				unit_link.participant_id,
				unit_link.customer_id,
				unit_link.investment_product_id,
				unit_link.ip,
				unit_link.transaction_date,
				unit_link.created_at,
				COALESCE(nav.amount, 0) as latest_nav,
				(unit_link.ip * COALESCE(nav.amount, 0)) as total_balance
			`).
			Joins("LEFT JOIN LATERAL (SELECT amount FROM net_asset_value WHERE investment_product_id = unit_link.investment_product_id AND deleted_at IS NULL AND created_date <= ? ORDER BY created_date DESC LIMIT 1) nav ON true", today).
			Where("unit_link.deleted_at is null")
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"participant_id"},
	})
	if err != nil {
		return nil, err
	}

	result := dataTable.(*pagination.ResultPagination)
	return result, nil
}
