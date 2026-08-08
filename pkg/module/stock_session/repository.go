package stocksession

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error)
	Get(ctx context.Context, id string) (*entity.StockSessionDto, error)
	GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error)
	Update(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	var result *entity.StockSessionDto
	err := r.db.Transaction(func(tx *gorm.DB) error {
		m := dto.ToModel()

		if err := tx.Omit("Items", "Payments", "Adjustments", "Logs").Create(m).Error; err != nil {
			return err
		}

		// Insert stock_session_item rows explicitly. We deliberately
		// do NOT rely on GORM's association auto-save here — that
		// would also try to INSERT into the item table whenever
		// `dto.Items[i].Item` is populated. The Omit("Item") is
		// belt-and-braces so even a stray nested Item cannot leak.
		// All we persist is the item_id reference + the row's own
		// snapshot fields.
		for _, it := range dto.Items {
			item := it.ToModel()
			item.SessionID = m.ID
			if err := tx.Omit("Item").Create(item).Error; err != nil {
				return err
			}
		}

		// Reload with associations so the call-site gets a fully
		// hydrated DTO back (items, payments, etc.).
		var reloaded model.StockSession
		if err := tx.
			Preload("Employee").
			Preload("Items").
			Preload("Items.Item").
			Where("id = ?", m.ID).First(&reloaded).Error; err != nil {
			return err
		}
		result = entity.NewStockSessionDtoFromModel(&reloaded)
		return nil
	})
	return result, err
}

func (r *repository) Get(ctx context.Context, id string) (*entity.StockSessionDto, error) {
	var m *model.StockSession
	if err := r.db.
		Preload("Employee").
		Preload("Items").
		Preload("Items.Item").
		Preload("Payments").
		Preload("Adjustments").
		Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewStockSessionDtoFromModel(m), nil
}

func (r *repository) GetByEmployeeDate(ctx context.Context, employeeID, date string) (*entity.StockSessionDto, error) {
	var m *model.StockSession
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	if err := r.db.
		Preload("Employee").
		Preload("Items").
		Preload("Items.Item").
		Where("employee_id = ? AND date = ?", employeeID, parsedDate).First(&m).Error; err != nil {
		return nil, err
	}
	return entity.NewStockSessionDtoFromModel(m), nil
}

func (r *repository) Update(ctx context.Context, dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	return r.updateInTx(dto)
}

func (r *repository) updateInTx(dto *entity.StockSessionDto) (*entity.StockSessionDto, error) {
	var result *entity.StockSessionDto
	err := r.db.Transaction(func(tx *gorm.DB) error {
		m := dto.ToModel()

		// Update the parent first via Save — this persists every
		// scalar column (totals, status, cash_debt, etc.). We do
		// NOT rely on association auto-save here because the
		// children slices below are handled explicitly. Auto-save
		// was producing inconsistent rows on some Go/GORM combos
		// (for example TotalQris landed as 0 because the children
		// failed silently on insert) — explicit writes are
		// predictable and the SQL log shows exactly what landed.
		if err := tx.Save(m).Error; err != nil {
			return err
		}

		// Wipe and re-insert child rows so the close path always
		// produces a clean row-set keyed on session_id. We
		// deliberately delete AFTER the parent Save to keep the
		// child-vs-parent ordering obvious in the SQL log.
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.StockSessionItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.PaymentDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", m.ID).Delete(&model.CashAdjustment{}).Error; err != nil {
			return err
		}

		// Insert stock_session_item rows explicitly with the
		// parent SessionID. Omit("Item") keeps any nested Item
		// DTO from leaking into an item-table INSERT.
		for _, it := range dto.Items {
			item := it.ToModel()
			item.SessionID = m.ID
			if err := tx.Omit("Item").Create(item).Error; err != nil {
				return err
			}
		}
		// Insert payment_detail rows explicitly. This is the row
		// that drives `total_qris`/`total_cash` — without the
		// explicit INSERT the totals computed by RecomputeTotals
		// had no payment rows backing them.
		for _, p := range dto.Payments {
			pay := p.ToModel()
			pay.SessionID = m.ID
			if err := tx.Create(pay).Error; err != nil {
				return err
			}
		}
		// Insert cash_adjustment rows explicitly.
		for _, a := range dto.Adjustments {
			adj := a.ToModel()
			adj.SessionID = m.ID
			if err := tx.Create(adj).Error; err != nil {
				return err
			}
		}

		// Reload with associations
		var reloaded model.StockSession
		if err := tx.
			Preload("Employee").
			Preload("Items").
			Preload("Items.Item").
			Preload("Payments").
			Preload("Adjustments").
			Where("id = ?", m.ID).First(&reloaded).Error; err != nil {
			return err
		}
		result = entity.NewStockSessionDtoFromModel(&reloaded)
		return nil
	})
	return result, err
}

// Delete removes a stock session row and all its child rows
// (stock_session_item / payment_detail / cash_adjustment) inside a
// single transaction so a partial failure cannot leave orphan
// children behind. Callers must have already verified the session is
// still OPEN — this repository method is a pure SQL primitive.
func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Child rows first so the parent delete can never race
		// with the cascading cleanup.
		if err := tx.Where("session_id = ?", id).Delete(&model.StockSessionItem{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", id).Delete(&model.PaymentDetail{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", id).Delete(&model.CashAdjustment{}).Error; err != nil {
			return err
		}
		// Parent last. Returns ErrRecordNotFound when no row
		// matched — the service layer maps that to a 404.
		if err := tx.Where("id = ?", id).Delete(&model.StockSession{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) FindAll(ctx context.Context, req *entity.StockSessionFindAllRequest) (*pagination.ResultPagination, error) {
	var m []model.StockSession = make([]model.StockSession, 0)
	tbl := pagination.NewTable(r.db)
	dataTable, err := tbl.Pagination(func(i interface{}) *gorm.DB {
		q := r.db.Model(&model.StockSession{})
		if req.FindAllRequest.OrganizationData.ID != "" {
			q = q.Where("organization_id = ?", req.FindAllRequest.OrganizationData.ID)
		}
		return q
	}, &pagination.TableRequest{
		Request:       req,
		QueryField:    []string{},
		Data:          &m,
		AllowedFields: []string{"date", "status", "total_sales"},
	})
	if err != nil {
		return nil, err
	}
	result := dataTable.(*pagination.ResultPagination)
	results := result.Data.(*[]model.StockSession)
	data := make([]*entity.StockSessionDto, 0)
	for i := range *results {
		// shallow hydrate (no nested) for list view — keeps payload small
		row := entity.NewStockSessionDtoFromModel(&(*results)[i])
		// Preload items via separate query (lightweight count only)
		data = append(data, row)
	}
	return &pagination.ResultPagination{
		Data:        data,
		Page:        result.Page,
		Count:       result.Count,
		RowsPerPage: result.RowsPerPage,
		TotalPages:  result.TotalPages,
	}, nil
}
