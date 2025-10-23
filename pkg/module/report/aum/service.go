package aum

import (
	"bytes"
	"context"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	GenerateReportAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]byte, error)
	GetReportAum(ctx context.Context, filter *entity.ReportAUMFilter) (*[]entity.GroupedAum, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

var indonesianMonths = [...]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func formatIndonesianDate(t time.Time) string {
	return fmt.Sprintf("%02d %s %d", t.Day(), indonesianMonths[t.Month()-1], t.Year())
}

func formatIndonesianMonthYear(month int, year int) string {
	if month >= 1 && month <= 12 {
		return fmt.Sprintf("%s %d", indonesianMonths[month-1], year)
	}
	return fmt.Sprintf("%d", year)
}

func formatCurrency(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func (s *service) GetReportAum(ctx context.Context, filter *entity.ReportAUMFilter) (*[]entity.GroupedAum, error) {

	companyData, err := s.repository.GetCompanyAUM(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get company AUM: %w", err)
	}

	var pesertaMandiriData []entity.ReportAUMData
	if filter.CompanyType == entity.ReportAUMCompanyTypePPIP {
		pesertaMandiriData, err = s.repository.GetPesertaMandiriAUM(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get peserta mandiri AUM: %w", err)
		}
	}

	allData := append(pesertaMandiriData, companyData...)

	var pilarData, nonPilarData []entity.ReportAUMData
	for _, data := range allData {
		if data.PilarType == "PILAR" {
			pilarData = append(pilarData, data)
		} else {
			nonPilarData = append(nonPilarData, data)
		}
	}

	result := []entity.GroupedAum{
		{
			Type:   entity.SummaryTypePilar,
			Groups: pilarData,
		},
		{
			Type:   entity.SummaryTypeNonPilar,
			Groups: nonPilarData,
		},
	}

	return &result, nil
}

func (s *service) GenerateReportAUM(ctx context.Context, filter *entity.ReportAUMFilter) ([]byte, error) {

	// Get company AUM data
	companyData, err := s.repository.GetCompanyAUM(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Get peserta mandiri data for PPIP
	var pesertaMandiriData []entity.ReportAUMData
	if filter.CompanyType == entity.ReportAUMCompanyTypePPIP {
		pesertaMandiriData, err = s.repository.GetPesertaMandiriAUM(ctx, filter)
		if err != nil {
			return nil, err
		}
	}

	allData := append(pesertaMandiriData, companyData...)

	// Create Excel file
	f := excelize.NewFile()
	sheet := "Report AUM"
	f.SetSheetName("Sheet1", sheet)

	// Set styles
	styleBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#00B0F0"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleBorder, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleBorderRight, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	styleBorderCenter, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Write header
	row := 1
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("Print Date : %s", formatIndonesianDate(time.Now())))
	row += 2

	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "SIM-55")
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBold)
	row++

	reportTitle := fmt.Sprintf("Report AUM %s", string(filter.CompanyType))
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), reportTitle)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBold)
	row++

	periode := fmt.Sprintf("Periode : %s", formatIndonesianMonthYear(filter.Month, filter.Year))
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), periode)
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBold)
	row += 2

	// Group data by pilar type
	pilarData := []entity.ReportAUMData{}
	nonPilarData := []entity.ReportAUMData{}

	for _, data := range allData {
		if data.PilarType == "PILAR" {
			pilarData = append(pilarData, data)
		} else {
			nonPilarData = append(nonPilarData, data)
		}
	}

	// Write PILAR section
	{
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "PILAR")
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBold)
		row++

		// Table headers
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "No")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Nama Perusahaan")
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Saldo Akhir (Rp)")
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), styleHeader)
		row++

		// Data rows
		var pilarTotal float64
		for i, data := range pilarData {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), data.CompanyName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatCurrency(data.TotalAUM))

			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBorderCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styleBorderRight)

			pilarTotal += data.TotalAUM
			row++
		}

		// Subtotal PILAR
		f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Subtotal Pilar")
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatCurrency(pilarTotal))

		subtotalStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})

		subtotalStyleRight, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true},
			Alignment: &excelize.Alignment{Horizontal: "right"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})

		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), subtotalStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), subtotalStyleRight)
		row += 2
	}

	// Write NON PILAR section
	{
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "NON PILAR")
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBold)
		row++

		// Table headers
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "No")
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Nama Perusahaan")
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Saldo Akhir (Rp)")
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("D%d", row), styleHeader)
		row++

		// Data rows
		var nonPilarTotal float64
		for i, data := range nonPilarData {
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), data.CompanyName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatCurrency(data.TotalAUM))

			f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), styleBorderCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styleBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styleBorderRight)

			nonPilarTotal += data.TotalAUM
			row++
		}

		// Subtotal NON PILAR
		f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Subtotal Non Pilar")
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatCurrency(nonPilarTotal))

		subtotalStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})

		subtotalStyleRight, _ := f.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Bold: true},
			Alignment: &excelize.Alignment{Horizontal: "right"},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
		})

		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), subtotalStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), subtotalStyleRight)
		row += 2
	}

	// Grand Total
	var grandTotal float64
	for _, data := range allData {
		grandTotal += data.TotalAUM
	}

	f.MergeCell(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row))
	grandTotalText := fmt.Sprintf("GRAND TOTAL %s", string(filter.CompanyType))
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), grandTotalText)
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), formatCurrency(grandTotal))

	grandTotalStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	grandTotalStyleRight, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("C%d", row), grandTotalStyle)
	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), grandTotalStyleRight)

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 3)
	f.SetColWidth(sheet, "B", "B", 6)
	f.SetColWidth(sheet, "C", "C", 30)
	f.SetColWidth(sheet, "D", "D", 25)

	// Set active sheet
	index, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(index)

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
