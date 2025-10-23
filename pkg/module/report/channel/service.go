package channel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	ei "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	netassetvalue "github.com/raymondsugiarto/coffee-api/pkg/module/net_asset_value"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	GenerateReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) ([]byte, error)
	GetTransactionReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) (*pagination.ResultPagination, error)
}

type service struct {
	repository           Repository
	netAssetValueService netassetvalue.Service
}

func NewService(repository Repository, netAssetValueService netassetvalue.Service) Service {
	return &service{
		repository:           repository,
		netAssetValueService: netAssetValueService,
	}
}

func formatIndonesianDate(t time.Time) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())

}

func (s *service) GetTransactionReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) (*pagination.ResultPagination, error) {
	adjustedEndDate := filter.EndDate.Add(7 * time.Hour)
	navItems, err := s.netAssetValueService.FindByDate(ctx, adjustedEndDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	nav := lo.KeyBy(navItems, func(item *ei.NetAssetValueDto) string {
		return item.InvestmentProductID
	})

	result, err := s.repository.GetTransactionReportChannel(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := result.Data.([]*entity.ReportTransactionChannel)

	for _, item := range items {
		if navItem, ok := nav[item.InvestmentProductID]; ok {
			item.Balance = float64(navItem.Amount) * item.Nab
		} else {
			return nil, errors.New("nav invalid")
		}
	}

	return result, nil
}

func (s *service) GenerateReportChannel(ctx context.Context, filter *entity.ReportTransactionChannelFilter) ([]byte, error) {
	adjustedEndDate := filter.EndDate.Add(7 * time.Hour)
	navItems, err := s.netAssetValueService.FindByDate(ctx, adjustedEndDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	nav := lo.KeyBy(navItems, func(item *ei.NetAssetValueDto) string {
		return item.InvestmentProductID
	})

	result, err := s.repository.GetTransactionReportChannel(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := result.Data.([]*entity.ReportTransactionChannel)

	excel := excelize.NewFile()
	sheetName := "Laporan Produksi"
	index, _ := excel.NewSheet(sheetName)

	boldStyle, _ := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	row := 1
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tanggal Cetak: "+formatIndonesianDate(time.Now()))
	row += 2

	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SIM-55")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "REPORT PRODUKSI")
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++
	excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("Periode: %s", formatIndonesianDate(adjustedEndDate)))
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row += 2

	// Header
	headers := []string{"NO", "NOMOR PESERTA", "NAMA PESERTA", "PRODUK", "NAMA PERUSAHAAN", "SALDO PESERTA", "AKUMULASI BIAYA"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		excel.SetCellValue(sheetName, cell, header)
	}

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
	excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("G%d", row), headerStyle)

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

	row++
	for i, item := range items {
		balance := 0.0
		if navItem, ok := nav[item.InvestmentProductID]; ok {
			balance = float64(navItem.Amount) * item.Nab
		} else {
			return nil, errors.New("nav invalid")
		}

		excel.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		excel.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.ParticipantID)
		excel.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.ParticipantName)
		excel.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.InvestmentProductName)
		excel.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.CompanyName)
		excel.SetCellValue(sheetName, fmt.Sprintf("F%d", row), balance)
		excel.SetCellValue(sheetName, fmt.Sprintf("G%d", row), item.Fee)
		excel.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("E%d", row), dataTextStyle)
		excel.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("G%d", row), numberStyle)

		row++
	}

	// Auto width
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		excel.SetColWidth(sheetName, colName, colName, 20)
	}
	// Set active sheet
	excel.SetActiveSheet(index)

	// Write to buffer
	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
