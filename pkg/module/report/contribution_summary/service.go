package contributionsummary

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
	GenerateContributionReport(ctx context.Context, filter *entity.ReportContributionSummaryFilter) ([]byte, error)
	GetContributionSummary(ctx context.Context, filter *entity.ReportContributionSummaryFilter) (*pagination.ResultPagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository}
}

func formatIndonesianDate(t time.Time) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func (s *service) GetContributionSummary(ctx context.Context, filter *entity.ReportContributionSummaryFilter) (*pagination.ResultPagination, error) {
	return s.repository.GetReportContribution(ctx, filter)
}

func (s *service) GenerateContributionReport(ctx context.Context, filter *entity.ReportContributionSummaryFilter) ([]byte, error) {
	result, err := s.repository.GetReportContribution(ctx, filter)
	if err != nil {
		return nil, err
	}

	data := result.Data.([]*entity.ReportContributionSummary)

	f := excelize.NewFile()
	sheetName := "Rincian Iuran"
	f.SetSheetName("Sheet1", sheetName)
	index, _ := f.NewSheet(sheetName)

	boldStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
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
	dataTextStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	numberStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt:    3,
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// ===== Header Info =====
	row := 1
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tanggal Cetak: "+formatIndonesianDate(time.Now()))
	row += 2

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SIM-55")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "REPORT RINCIAN IURAN PPIP")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("Periode: %s", formatIndonesianDate(*filter.EndDate)))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row += 2

	// ===== Table Header =====
	headers := []string{
		"No", "Nama Kelompok/Individu", "Iuran Normal Peserta",
		"Iuran Sukarela Pekerja", "Iuran Normal Pemberi Kerja", "Iuran Dana Pendidikan Anak Peserta",
		"Total", "Program Pensiun/Manfaat Lain",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		f.SetCellValue(sheetName, cell, h)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("H%d", row), headerStyle)
	row++

	// ===== Data Rows =====
	startRow := row
	for i, col := range data {
		r := startRow + i

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", r), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", r), col.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", r), col.CustomerAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", r), col.VoluntaryAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", r), col.EmployerAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", r), col.EducationFundAmount)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", r), col.Total)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", r), col.TypeCode)

		// Apply styles
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", r), fmt.Sprintf("B%d", r), dataTextStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", r), fmt.Sprintf("G%d", r), numberStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("H%d", r), fmt.Sprintf("H%d", r), dataTextStyle)
	}

	// ===== Footer Notes =====
	footerStart := startRow + len(data) + 1
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", footerStart), "1001 = PPIP")
	f.SetCellStyle(sheetName, fmt.Sprintf("H%d", footerStart), fmt.Sprintf("G%d", footerStart), boldStyle)

	f.SetCellValue(sheetName, fmt.Sprintf("H%d", footerStart+1), "1002 = DKP dan Manfaat Lain")
	f.SetCellStyle(sheetName, fmt.Sprintf("H%d", footerStart+1), fmt.Sprintf("G%d", footerStart+1), boldStyle)

	// ===== Auto Column Width =====
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 22)
	}

	f.SetActiveSheet(index)

	// ===== Output to buffer =====
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
