package participantsummary

import (
	"bytes"
	"context"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/module/benefit_type"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	GenerateReportParticipantSummary(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]byte, error)
	GetParticipantSummary(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.GroupedSummary, error)
}

type service struct {
	repository         Repository
	benefitTypeService benefit_type.Service
}

func NewService(repository Repository, benefitTypeService benefit_type.Service) Service {
	return &service{
		repository:         repository,
		benefitTypeService: benefitTypeService,
	}
}

func formatIndonesianDate(t time.Time) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func formatNumber(n int64) string {
	return fmt.Sprintf("%d", n)
}

var defaultCombinations = []entity.ReportParticipantSummary{
	{ParticipantType: "PPIP", ParticipantCategory: "MANDIRI", PilarType: "ALL"},
	{ParticipantType: "PPIP", ParticipantCategory: "KUMPULAN", PilarType: "ALL"},
	{ParticipantType: "DKP", ParticipantCategory: "KUMPULAN", PilarType: "ALL"},

	{ParticipantType: "PPIP", ParticipantCategory: "MANDIRI", PilarType: "PILAR"},
	{ParticipantType: "PPIP", ParticipantCategory: "KUMPULAN", PilarType: "PILAR"},
	{ParticipantType: "DKP", ParticipantCategory: "KUMPULAN", PilarType: "PILAR"},

	{ParticipantType: "PPIP", ParticipantCategory: "MANDIRI", PilarType: "NON PILAR"},
	{ParticipantType: "PPIP", ParticipantCategory: "KUMPULAN", PilarType: "NON PILAR"},
	{ParticipantType: "DKP", ParticipantCategory: "KUMPULAN", PilarType: "NON PILAR"},
}

func filterExpected(pilarType string) []entity.ReportParticipantSummary {
	var list []entity.ReportParticipantSummary
	for _, c := range defaultCombinations {
		if c.PilarType == pilarType {
			list = append(list, c)
		}
	}
	return list
}

func padSummary(data []entity.ReportParticipantSummary, expected []entity.ReportParticipantSummary) []entity.ReportParticipantSummary {
	result := make([]entity.ReportParticipantSummary, 0, len(expected))
	for _, ex := range expected {
		found := false
		for _, row := range data {
			if row.ParticipantType == ex.ParticipantType &&
				row.ParticipantCategory == ex.ParticipantCategory &&
				row.PilarType == ex.PilarType {
				result = append(result, row)
				found = true
				break
			}
		}
		if !found {
			ex.ParticipantCount = 0
			result = append(result, ex)
		}
	}
	return result
}

func padBenefitSummary(data []entity.ReportParticipantSummary, benefitTypes []string) []entity.ReportParticipantSummary {
	var result []entity.ReportParticipantSummary
	categories := []string{"MANDIRI", "KUMPULAN"}

	for _, bt := range benefitTypes {
		for _, cat := range categories {
			found := false
			for _, d := range data {
				if d.ParticipantType == bt && d.ParticipantCategory == cat {
					result = append(result, d)
					found = true
					break
				}
			}
			if !found {
				result = append(result, entity.ReportParticipantSummary{
					ParticipantType:     bt,
					ParticipantCategory: cat,
					ParticipantCount:    0,
					PilarType:           "",
				})
			}
		}
	}
	return result
}

func padBusinessGroupSummary(data []entity.ReportParticipantSummary) []entity.ReportParticipantSummary {
	dataMap := map[string]entity.ReportParticipantSummary{}
	for _, row := range data {
		if row.DataSource == "PERUSAHAAN" {
			key := row.ParticipantType + "_" + row.ParticipantCategory
			dataMap[key] = row
		}
	}

	orderedCombinations := []struct {
		companyType string
		pilarType   string
	}{
		{"DKP", "PILAR"},
		{"DKP", "NON PILAR"},
		{"PPIP", "PILAR"},
		{"PPIP", "NON PILAR"},
		{"DKP", ""},
		{"PPIP", ""},
	}

	result := []entity.ReportParticipantSummary{}
	for _, combo := range orderedCombinations {
		key := combo.companyType + "_" + combo.pilarType
		if row, exists := dataMap[key]; exists {
			result = append(result, row)
		} else {
			result = append(result, entity.ReportParticipantSummary{
				ParticipantType:     combo.companyType,
				ParticipantCategory: combo.pilarType,
				ParticipantCount:    0,
				PilarType:           "",
				DataSource:          "PERUSAHAAN",
			})
		}
	}

	return result
}

func (s *service) GetParticipantSummary(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]entity.GroupedSummary, error) {
	all, err := s.repository.GetSummaryAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	pilar, err := s.repository.GetSummaryPilar(ctx, filter)
	if err != nil {
		return nil, err
	}
	nonPilar, err := s.repository.GetSummaryNonPilar(ctx, filter)
	if err != nil {
		return nil, err
	}
	manfaatLain, err := s.repository.GetSummaryManfaatLain(ctx, filter)
	if err != nil {
		return nil, err
	}
	perusahaan, err := s.repository.GetSummaryPerusahaan(ctx, filter)
	if err != nil {
		return nil, err
	}

	// === Normalize biar sama dengan GenerateReportParticipantSummary ===
	normalizedAll := padSummary(all, filterExpected("ALL"))
	normalizedPilar := padSummary(pilar, filterExpected("PILAR"))
	normalizedNonPilar := padSummary(nonPilar, filterExpected("NON PILAR"))
	normalizedPerusahaan := padBusinessGroupSummary(perusahaan)

	benefitTypeNames, err := s.benefitTypeService.GetActiveBenefitTypeNames(ctx)
	if err != nil {
		return nil, err
	}
	normalizedManfaatLain := padBenefitSummary(manfaatLain, benefitTypeNames)

	result := []entity.GroupedSummary{
		{
			Type:   entity.SummaryTypeAll,
			Groups: normalizedAll,
		},
		{
			Type:   entity.SummaryTypePilar,
			Groups: normalizedPilar,
		},
		{
			Type:   entity.SummaryTypeNonPilar,
			Groups: normalizedNonPilar,
		},
		{
			Type:   entity.SummaryTypeManfaatLain,
			Groups: normalizedManfaatLain,
		},
		{
			Type:   entity.SummaryTypePesertaBadanUsaha,
			Groups: normalizedPerusahaan,
		},
	}

	return result, nil
}

func (s *service) GenerateReportParticipantSummary(ctx context.Context, filter *entity.ReportParticipantSummaryFilter) ([]byte, error) {
	all, _ := s.repository.GetSummaryAll(ctx, filter)
	pilar, _ := s.repository.GetSummaryPilar(ctx, filter)
	nonPilar, _ := s.repository.GetSummaryNonPilar(ctx, filter)
	manfaatLain, _ := s.repository.GetSummaryManfaatLain(ctx, filter)
	perusahaan, _ := s.repository.GetSummaryPerusahaan(ctx, filter)

	normalizedAll := padSummary(all, filterExpected("ALL"))
	normalizedPilar := padSummary(pilar, filterExpected("PILAR"))
	normalizedNonPilar := padSummary(nonPilar, filterExpected("NON PILAR"))
	normalizedPerusahaan := padBusinessGroupSummary(perusahaan)

	benefitTypeNames, err := s.benefitTypeService.GetActiveBenefitTypeNames(ctx)
	if err != nil {
		return nil, err
	}
	normalizedManfaatLain := padBenefitSummary(manfaatLain, benefitTypeNames)

	excel := excelize.NewFile()
	sheetName := "Summary Peserta"
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

	// Header
	row := 1
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tanggal Cetak: "+formatIndonesianDate(time.Now()))
	row += 2

	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SIM-55")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Summary Peserta")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Periode: "+formatIndonesianDate(*filter.EndDate))
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row += 2

	// Tulis tiap section
	writeSectionWithStyle(excel, sheetName, &row, "ALL", normalizedAll, headerStyle, dataTextStyle, numberStyle)
	writeSectionWithStyle(excel, sheetName, &row, "PILAR", normalizedPilar, headerStyle, dataTextStyle, numberStyle)
	writeSectionWithStyle(excel, sheetName, &row, "NON PILAR", normalizedNonPilar, headerStyle, dataTextStyle, numberStyle)
	writeSectionWithStyle(excel, sheetName, &row, "MANFAAT LAIN", normalizedManfaatLain, headerStyle, dataTextStyle, numberStyle)
	writeSectionWithStyle(excel, sheetName, &row, "PESERTA BADAN USAHA", normalizedPerusahaan, headerStyle, dataTextStyle, numberStyle)

	// Auto width
	for col := 1; col <= 4; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		excel.SetColWidth(sheetName, colName, colName, 20)
	}

	// Set active sheet
	excel.SetActiveSheet(index)

	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSectionWithStyle(f *excelize.File, sheet string, row *int, title string, data []entity.ReportParticipantSummary, headerStyle, dataTextStyle, numberStyle int) {

	sectionTitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), title)
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("A%d", *row), sectionTitleStyle)
	*row++

	headers := []string{"No", "PRODUK", "KATEGORI", "JUMLAH PESERTA"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, *row)
		f.SetCellValue(sheet, cell, h)
	}
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("D%d", *row), headerStyle)
	*row++

	var total int64
	for i, rowData := range data {
		no := i + 1
		produk := rowData.ParticipantType
		kategori := rowData.ParticipantCategory

		f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), no)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", *row), produk)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", *row), kategori)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", *row), rowData.ParticipantCount)

		f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("C%d", *row), dataTextStyle)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", *row), fmt.Sprintf("D%d", *row), numberStyle)

		total += rowData.ParticipantCount
		*row++
	}

	grandTotalStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	grandTotalNumberStyle, _ := f.NewStyle(&excelize.Style{
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

	f.MergeCell(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("C%d", *row))
	f.SetCellValue(sheet, fmt.Sprintf("A%d", *row), "GRAND TOTAL")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", *row), total)

	f.SetCellStyle(sheet, fmt.Sprintf("A%d", *row), fmt.Sprintf("C%d", *row), grandTotalStyle)
	f.SetCellStyle(sheet, fmt.Sprintf("D%d", *row), fmt.Sprintf("D%d", *row), grandTotalNumberStyle)

	*row += 2
}
