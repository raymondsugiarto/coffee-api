package stocksession

import (
	"context"
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"gorm.io/gorm"
)

func (s *service) GetDashboard(ctx context.Context) (*entity.DashboardSummaryDto, error) {
	today := time.Now().Format("2006-01-02")

	type aggRow struct {
		TotalSales     float64
		TotalCash      float64
		TotalQris      float64
		Transactions   int64
		OpenSessions   int64
		ClosedSessions int64
		TotalSessions  int64
	}

	var row aggRow
	err := s.db.Raw(`
		SELECT
			COALESCE(SUM(total_sales), 0) AS total_sales,
			COALESCE(SUM(total_cash), 0) AS total_cash,
			COALESCE(SUM(total_qris), 0) AS total_qris,
			COUNT(*) AS transactions,
			SUM(CASE WHEN status = 'OPEN' THEN 1 ELSE 0 END) AS open_sessions,
			SUM(CASE WHEN status = 'CLOSED' THEN 1 ELSE 0 END) AS closed_sessions,
			COUNT(*) AS total_sessions
		FROM stock_session
		WHERE date = ?
	`, today).Scan(&row).Error
	if err != nil {
		return nil, err
	}

	// Fallback: also sum payments directly for today's transactions (covers OPEN sessions)
	type payAgg struct {
		TotalCash float64
		TotalQris float64
	}
	var p payAgg
	s.db.Raw(`
		SELECT
			COALESCE(SUM(CASE WHEN pd.payment_method = 'CASH' THEN pd.amount ELSE 0 END), 0) AS total_cash,
			COALESCE(SUM(CASE WHEN pd.payment_method = 'QRIS' THEN pd.amount ELSE 0 END), 0) AS total_qris
		FROM payment_detail pd
		JOIN stock_session ss ON ss.id = pd.session_id
		WHERE ss.date = ?
	`, today).Scan(&p)

	return &entity.DashboardSummaryDto{
		TodaySales:        row.TotalSales,
		TodayCash:         row.TotalCash + p.TotalCash,
		TodayQris:         row.TotalQris + p.TotalQris,
		TodayTransactions: int(row.Transactions),
		OpenSessions:      int(row.OpenSessions),
		ClosedSessions:    int(row.ClosedSessions),
		TotalSessions:     int(row.TotalSessions),
	}, nil
}

func (s *service) GetDailyReport(ctx context.Context, date string) (*entity.DailyReportDto, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	report := &entity.DailyReportDto{Date: date}

	var sessions []model.StockSession
	if err := s.db.Where("date = ?", date).Find(&sessions).Error; err != nil {
		return nil, err
	}

	byEmployee := make(map[string]*entity.EmployeeReportRowDto)
	for _, ss := range sessions {
		report.Sessions++
		report.TotalSales += ss.TotalSales
		report.TotalCash += ss.TotalCash
		report.TotalQris += ss.TotalQris
		report.TotalOther += ss.TotalOther
		report.TotalPayment += ss.TotalPayment
		report.TotalDiff += ss.Difference
		report.TotalCommission += ss.TotalCommission
		report.TotalMealAllowance += ss.MealAllowance
		report.TotalBonusTarget += ss.BonusTarget
		report.TotalSalary += ss.TotalSalary

		row, ok := byEmployee[ss.EmployeeID]
		if !ok {
			row = &entity.EmployeeReportRowDto{EmployeeID: ss.EmployeeID}
			byEmployee[ss.EmployeeID] = row
		}
		row.Sessions++
		row.TotalSales += ss.TotalSales
		row.TotalCash += ss.TotalCash
		row.TotalQris += ss.TotalQris
		row.Difference += ss.Difference
		row.Commission += ss.TotalCommission
		row.MealAllowance += ss.MealAllowance
		row.BonusTarget += ss.BonusTarget
		row.TotalSalary += ss.TotalSalary
	}

	// Hydrate employee names
	if len(byEmployee) > 0 {
		ids := make([]string, 0, len(byEmployee))
		for id := range byEmployee {
			ids = append(ids, id)
		}
		var drivers []model.Admin
		if err := s.db.Where("id IN ?", ids).Find(&drivers).Error; err == nil {
			for _, d := range drivers {
				if r, ok := byEmployee[d.ID]; ok {
					r.EmployeeName = (d.FirstName + " " + d.LastName)
				}
			}
		}
		for _, r := range byEmployee {
			report.ByEmployee = append(report.ByEmployee, *r)
		}
	}
	return report, nil
}

func (s *service) GetMonthlyReport(ctx context.Context, year, month int) (*entity.MonthlyReportDto, error) {
	if year == 0 {
		year = time.Now().Year()
	}
	if month == 0 {
		month = int(time.Now().Month())
	}
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)

	report := &entity.MonthlyReportDto{Year: year, Month: month}

	var sessions []model.StockSession
	if err := s.db.Where("date >= ? AND date < ?", from, to).Find(&sessions).Error; err != nil {
		return nil, err
	}

	byDate := make(map[string]*entity.DailyReportDto)
	byEmployee := make(map[string]*entity.EmployeeReportRowDto)
	for _, ss := range sessions {
		dateKey := ss.Date.Format("2006-01-02")
		report.Sessions++
		report.TotalSales += ss.TotalSales
		report.TotalCash += ss.TotalCash
		report.TotalQris += ss.TotalQris
		report.TotalDiff += ss.Difference
		report.TotalCommission += ss.TotalCommission
		report.TotalMealAllowance += ss.MealAllowance
		report.TotalBonusTarget += ss.BonusTarget
		report.TotalSalary += ss.TotalSalary

		daily, ok := byDate[dateKey]
		if !ok {
			daily = &entity.DailyReportDto{Date: dateKey}
			byDate[dateKey] = daily
		}
		daily.Sessions++
		daily.TotalSales += ss.TotalSales
		daily.TotalCash += ss.TotalCash
		daily.TotalQris += ss.TotalQris
		daily.TotalDiff += ss.Difference
		daily.TotalCommission += ss.TotalCommission
		daily.TotalMealAllowance += ss.MealAllowance
		daily.TotalBonusTarget += ss.BonusTarget
		daily.TotalSalary += ss.TotalSalary

		row, ok := byEmployee[ss.EmployeeID]
		if !ok {
			row = &entity.EmployeeReportRowDto{EmployeeID: ss.EmployeeID}
			byEmployee[ss.EmployeeID] = row
		}
		row.Sessions++
		row.TotalSales += ss.TotalSales
		row.TotalCash += ss.TotalCash
		row.TotalQris += ss.TotalQris
		row.Difference += ss.Difference
		row.Commission += ss.TotalCommission
		row.MealAllowance += ss.MealAllowance
		row.BonusTarget += ss.BonusTarget
		row.TotalSalary += ss.TotalSalary
	}

	if len(byDate) > 0 {
		for _, d := range byDate {
			report.Daily = append(report.Daily, *d)
		}
	}
	if len(byEmployee) > 0 {
		ids := make([]string, 0, len(byEmployee))
		for id := range byEmployee {
			ids = append(ids, id)
		}
		var drivers []model.Admin
		if err := s.db.Where("id IN ?", ids).Find(&drivers).Error; err == nil {
			for _, d := range drivers {
				if r, ok := byEmployee[d.ID]; ok {
					r.EmployeeName = (d.FirstName + " " + d.LastName)
				}
			}
		}
		for _, r := range byEmployee {
			report.ByEmployee = append(report.ByEmployee, *r)
		}
	}
	return report, nil
}

func (s *service) GetTopProducts(ctx context.Context, from, to string, limit int) ([]entity.TopProductRowDto, error) {
	if limit <= 0 {
		limit = 10
	}
	if from == "" {
		from = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}

	type rawRow struct {
		ProductID   string
		ProductName string
		SKU         string
		TotalQty    int
		TotalSales  float64
	}
	var rows []rawRow
	err := s.db.
		Table("stock_session_item ssi").
		Select("ssi.item_id as item_id, p.name as product_name, COALESCE(p.sku, p.code) as sku, SUM(ssi.sold_qty) as total_qty, SUM(ssi.subtotal) as total_sales").
		Joins("JOIN stock_session ss ON ss.id = ssi.session_id").
		Joins("JOIN item p ON p.id = ssi.item_id").
		Where("ss.date >= ? AND ss.date <= ?", from, to).
		Group("ssi.item_id, p.name, p.sku").
		Order("total_qty DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.TopProductRowDto, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.TopProductRowDto{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			SKU:         r.SKU,
			TotalQty:    r.TotalQty,
			TotalSales:  r.TotalSales,
		})
	}
	return out, nil
}

func (s *service) GetEmployeePerformance(ctx context.Context, from, to string) ([]entity.EmployeePerformanceRowDto, error) {
	if from == "" {
		from = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}
	type rawRow struct {
		EmployeeID    string
		FirstName     string
		LastName      string
		Sessions      int64
		TotalItems    int64
		TotalSales    float64
		TotalCash     float64
		TotalQris     float64
		TotalDiff     float64
		Commission    float64
		MealAllowance float64
		BonusTarget   float64
		TotalSalary   float64
	}
	var rows []rawRow
	err := s.db.
		Table("stock_session ss").
		Select(`ss.employee_id as employee_id,
		        a.first_name as first_name,
		        a.last_name as last_name,
		        COUNT(*) as sessions,
		        COALESCE(SUM(ss.total_items), 0) as total_items,
		        COALESCE(SUM(ss.total_sales), 0) as total_sales,
		        COALESCE(SUM(ss.total_cash), 0) as total_cash,
		        COALESCE(SUM(ss.total_qris), 0) as total_qris,
		        COALESCE(SUM(ss.difference), 0) as total_diff,
		        COALESCE(SUM(ss.total_commission), 0) as commission,
		        COALESCE(SUM(ss.meal_allowance), 0) as meal_allowance,
		        COALESCE(SUM(ss.bonus_target), 0) as bonus_target,
		        COALESCE(SUM(ss.total_salary), 0) as total_salary`).
		Joins("LEFT JOIN admin a ON a.id = ss.employee_id").
		Where("ss.date >= ? AND ss.date <= ?", from, to).
		Group("ss.employee_id, a.first_name, a.last_name").
		Order("total_sales DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]entity.EmployeePerformanceRowDto, 0, len(rows))
	for _, r := range rows {
		out = append(out, entity.EmployeePerformanceRowDto{
			EmployeeID:    r.EmployeeID,
			EmployeeName:  (r.FirstName + " " + r.LastName),
			Sessions:      int(r.Sessions),
			TotalItems:    int(r.TotalItems),
			TotalSales:    r.TotalSales,
			TotalCash:     r.TotalCash,
			TotalQris:     r.TotalQris,
			TotalDiff:     r.TotalDiff,
			Commission:    r.Commission,
			MealAllowance: r.MealAllowance,
			BonusTarget:   r.BonusTarget,
			TotalSalary:   r.TotalSalary,
		})
	}
	return out, nil
}

// ensure gorm.DB is referenced for compilation
var _ = gorm.ErrRecordNotFound
