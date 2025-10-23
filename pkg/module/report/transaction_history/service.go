package transactionhistory

import (
	"context"
	"fmt"
	"time"

	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	GetTransactionHistoryReport(ctx context.Context, filter *entity.TransactionHistoryFilter) (*pagination.ResultPagination, error)
	GenerateExcel(ctx context.Context, req *entity.TransactionHistoryFilter, result *pagination.ResultPagination) ([]byte, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func formatIndonesianDate(t time.Time) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())
}

func (s *service) GetTransactionHistoryReport(ctx context.Context, req *entity.TransactionHistoryFilter) (*pagination.ResultPagination, error) {
	req.Status = model.InvestmentStatusSuccess
	return s.repo.GetTransactionHistory(ctx, req)
}

func (s *service) GenerateExcel(ctx context.Context, req *entity.TransactionHistoryFilter, result *pagination.ResultPagination) ([]byte, error) {
	req.Status = model.InvestmentStatusSuccess

	f := excelize.NewFile()
	sheetName := "Riwayat Transaksi"
	index, _ := f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

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

	var periode time.Time
	if !req.EndDate.IsZero() {
		periode = req.EndDate
	} else if !req.StartDate.IsZero() {
		periode = req.StartDate
	} else {
		periode = time.Now()
	}

	row := 1
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "Tanggal Cetak: "+formatIndonesianDate(time.Now()))
	row += 2

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "SIM-55")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "REPORT RIWAYAT TRANSAKSI")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("Periode: %s", formatIndonesianDate(periode)))
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), boldStyle)
	row += 2

	headers := []string{
		"Kode Transaksi", "Tgl Transaksi", "Kode Peserta", "Nama Peserta",
		"Produk Investasi", "Nominal Transaksi", "NAB Transaksi", "Unit Transaksi",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, row)
		f.SetCellValue(sheetName, cell, header)
	}
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("H%d", row), headerStyle)
	row++

	startRow := row
	data := result.Data.([]*entity.TransactionHistoryReport)
	for i, item := range data {
		row := startRow + i

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.InvestmentItemCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.InvestmentAt.Add(7*time.Hour).Format("2006-01-02 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.ParticipantCode)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.ParticipantName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.InvestmentProductName)

		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), item.TotalAmount)

		if item.NavAmount != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), *item.NavAmount)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "-")
		}

		if item.UnitAmount > 0 {
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), item.UnitAmount)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), "-")
		}

		// Apply styles
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("E%d", row), dataTextStyle)
		f.SetCellStyle(sheetName, fmt.Sprintf("F%d", row), fmt.Sprintf("H%d", row), numberStyle)
	}

	// ===== Auto Column Width =====
	for col := 1; col <= len(headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 20)
	}

	f.SetActiveSheet(index)

	// ===== Output to buffer =====
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
