package participantsummary

import (
	"bytes"
	"context"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	GenerateReportSummaryAum(ctx context.Context, filter *entity.ReportSummaryAumFilter) ([]byte, error)
	GetSummaryAum(ctx context.Context, filter *entity.ReportSummaryAumFilter) (*pagination.ResultPagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func formatIndonesianDate(t time.Time) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func indonesianMonthName(month int) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	return months[month-1]
}

func (s *service) GetSummaryAum(ctx context.Context, filter *entity.ReportSummaryAumFilter) (*pagination.ResultPagination, error) {
	raw, err := s.repository.GetTransactionPerCompanyTypeAndPilar(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := raw.Data.([]*entity.ReportSummaryAum)

	// Agregasi
	pilar := &entity.ReportSummaryAumAggregated{GroupProduksi: "PILAR"}
	nonPilar := &entity.ReportSummaryAumAggregated{GroupProduksi: "NON PILAR"}

	for _, summary := range items {
		if summary.CompanyType == "PPIP" && summary.PilarType == "PILAR" {
			pilar.AumPPIP += summary.Aum
		} else if summary.CompanyType == "PPIP" {
			nonPilar.AumPPIP += summary.Aum
		} else if summary.CompanyType == "DKP" && summary.PilarType == "PILAR" {
			pilar.AumDKP += summary.Aum
		} else {
			nonPilar.AumDKP += summary.Aum
		}
	}

	pilar.TotalAum = pilar.AumPPIP + pilar.AumDKP
	nonPilar.TotalAum = nonPilar.AumPPIP + nonPilar.AumDKP

	aggregated := []*entity.ReportSummaryAumAggregated{pilar, nonPilar}

	return &pagination.ResultPagination{
		Data:        aggregated,
		Page:        raw.Page,
		Count:       raw.Count,
		RowsPerPage: raw.RowsPerPage,
		TotalPages:  raw.TotalPages,
	}, nil
}

func (s *service) GenerateReportSummaryAum(ctx context.Context, filter *entity.ReportSummaryAumFilter) ([]byte, error) {
	summaryAum, err := s.repository.GetTransactionPerCompanyTypeAndPilar(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := summaryAum.Data.([]*entity.ReportSummaryAum)

	pilarAumDKP := 0.0
	pilarAumPPIP := 0.0
	nonPilarAumPPIP := 0.0
	nonPilarAumDKP := 0.0

	for _, summary := range items {
		if summary.CompanyType == "PPIP" && summary.PilarType == "PILAR" {
			pilarAumPPIP += summary.Aum
		} else if summary.CompanyType == "PPIP" {
			nonPilarAumPPIP += summary.Aum
		} else if summary.CompanyType == "DKP" && summary.PilarType == "PILAR" {
			pilarAumDKP += summary.Aum
		} else {
			nonPilarAumDKP += summary.Aum
		}
	}

	excel := excelize.NewFile()
	sheetName := "Summary AUM"
	excel.SetSheetName("Sheet1", sheetName)
	index, _ := excel.NewSheet(sheetName)

	// Styles
	boldStyle, _ := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	headerStyle, _ := excel.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	dataTextStyle, _ := excel.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	numberStyle, _ := excel.NewStyle(&excelize.Style{
		NumFmt:    3,
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	grandTotalStyle, _ := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	grandTotalNumberStyle, _ := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt:    3,
		Alignment: &excelize.Alignment{Horizontal: "right"},
	})

	// Header Info
	row := 1
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tanggal Cetak: "+formatIndonesianDate(time.Now()))
	row += 2

	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SIM-55")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SUMMARY AUM")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("Periode: %s %d", indonesianMonthName(int(filter.Month)), filter.Year))
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row += 2
	// Table Header
	headers := []string{"NO", "GROUP PRODUKSI", "AUM PPIP", "AUM DKP", "TOTAL AUM"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		excel.SetCellValue(sheetName, cell, h)
	}
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("E%d", row), headerStyle)
	row++

	// Data Rows
	data := []struct {
		GroupProduksi string
		AumPPIP       float64
		AumDKP        float64
	}{
		{"PILAR", pilarAumPPIP, pilarAumDKP},
		{"NON PILAR", nonPilarAumPPIP, nonPilarAumDKP},
	}

	var totalPPIP, totalDKP float64
	for i, d := range data {
		total := d.AumPPIP + d.AumDKP

		excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		excel.SetCellValue(sheetName, fmt.Sprintf("B%d", row), d.GroupProduksi)
		excel.SetCellValue(sheetName, fmt.Sprintf("C%d", row), d.AumPPIP)
		excel.SetCellValue(sheetName, fmt.Sprintf("D%d", row), d.AumDKP)
		excel.SetCellValue(sheetName, fmt.Sprintf("E%d", row), total)

		excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), dataTextStyle)
		excel.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("E%d", row), numberStyle)

		totalPPIP += d.AumPPIP
		totalDKP += d.AumDKP
		row++
	}

	// Grand Total
	totalAUM := totalPPIP + totalDKP
	excel.MergeCell(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row))
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "GRAND TOTAL")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), grandTotalStyle)

	excel.SetCellValue(sheetName, fmt.Sprintf("C%d", row), totalPPIP)
	excel.SetCellValue(sheetName, fmt.Sprintf("D%d", row), totalDKP)
	excel.SetCellValue(sheetName, fmt.Sprintf("E%d", row), totalAUM)
	excel.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("E%d", row), numberStyle)
	excel.SetCellStyle(sheetName, fmt.Sprintf("C%d", row), fmt.Sprintf("E%d", row), grandTotalNumberStyle)

	// Auto width
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		excel.SetColWidth(sheetName, colName, colName, 20)
	}

	excel.SetActiveSheet(index)

	// Buffer Output
	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
