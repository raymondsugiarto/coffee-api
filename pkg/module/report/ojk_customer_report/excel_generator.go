package ojkcustomerreport

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer_report"
	"github.com/xuri/excelize/v2"
)

// ExcelStyles contains all Excel cell styles for the report
type ExcelStyles struct {
	styleHeader                  int
	styleBold                    int
	styleBorder                  int
	styleBorderRight             int
	styleBorderCenter            int
	styleNumber                  int
	styleNumberBorder            int
	styleNumberBorderRight       int
	styleSubtotal                int
	styleSubtotalRight           int
	styleGrandTotal              int
	styleGrandTotalRight         int
	styleGrayBackground          int
	styleNoBorder                int
	styleAlternateRow            int
	styleAlternateRowCenter      int
	styleAlternateRowNumber      int
	styleAlternateRowNumberRight int
	styleCurrencyBorder          int
	styleCurrencyAlternateRow    int
	styleCurrencyHeader          int
	styleCurrencyNormal          int
	styleCurrencyBorderBold      int
	styleBoldRight               int
	styleInfoBorder              int
	styleInfoBorderTop           int
	styleInfoBorderMiddle        int
	styleInfoBorderBottom        int
	styleBottomBorderOnly        int
	styleBottomBorderBold        int
	styleBottomBorderRight       int
	styleBottomBorderNumber      int
	styleBottomBorderCurrency    int
	styleNormal                  int
}

type ExcelGenerator struct {
}

func NewExcelGenerator() *ExcelGenerator {
	return &ExcelGenerator{}
}

func (g *ExcelGenerator) GenerateExcel(data *crEntity.OJKCustomerReportDataDto) ([]byte, error) {
	f := excelize.NewFile()
	sheet := "Customer Report"
	f.SetSheetName("Sheet1", sheet)

	// Create all Excel styles
	styles := g.createComprehensiveStyles(f)

	row := 3 // Start at row 3 for border padding

	// Write main header
	row = g.writeMainHeader(f, sheet, row, data, styles)
	row += 2

	// Write customer information
	row = g.writeCustomerInfoSection(f, sheet, row, data, styles)
	row += 2

	// Write transaction sections
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "• RINGKASAN TRANSAKSI")
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleBold)
	row += 1

	for _, section := range data.TransactionSections {
		row = g.writeTransactionSection(f, sheet, row, section, styles)
		row += 2
	}

	// Write side-by-side sections
	row = g.writeSideBySideSections(f, sheet, row, data, styles)
	row += 2

	// Write footer information
	row = g.writeFooterSection(f, sheet, row, styles)

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 3)
	f.SetColWidth(sheet, "B", "B", 2)  // Left padding border
	f.SetColWidth(sheet, "C", "C", 30) // Jenis Transaksi
	f.SetColWidth(sheet, "D", "D", 20) // Tanggal Transaksi
	f.SetColWidth(sheet, "E", "E", 20) // Nilai Investasi (a)
	f.SetColWidth(sheet, "F", "F", 25) // Tanggal Harga Unit
	f.SetColWidth(sheet, "G", "G", 20) // Harga per Unit (b)
	f.SetColWidth(sheet, "H", "H", 20) // Transaksi a/b (unit)
	f.SetColWidth(sheet, "I", "I", 25) // Saldo (unit)
	f.SetColWidth(sheet, "J", "J", 2)  // Right padding border

	// Apply white background to empty areas
	whiteBackgroundStyle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Fill left margin
	for i := 1; i <= row+10; i++ {
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), whiteBackgroundStyle)
	}

	// Fill right areas
	for i := 1; i <= row+10; i++ {
		for col := 'K'; col <= 'Z'; col++ {
			f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, i), fmt.Sprintf("%c%d", col, i), whiteBackgroundStyle)
		}
	}

	// Fill empty cells with white background
	for i := 1; i <= row; i++ {
		for col := 'C'; col <= 'I'; col++ {
			cellRef := fmt.Sprintf("%c%d", col, i)
			cellValue, _ := f.GetCellValue(sheet, cellRef)
			// Apply white background to empty cells only
			if cellValue == "" {
				f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, i), fmt.Sprintf("%c%d", col, i), whiteBackgroundStyle)
			}
		}
	}

	// Fill rows below content
	for i := row + 1; i <= row+20; i++ {
		for col := 'A'; col <= 'Z'; col++ {
			f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, i), fmt.Sprintf("%c%d", col, i), whiteBackgroundStyle)
		}
	}

	// Add report border frame (after white background to preserve borders)
	g.addReportBorderFrame(f, sheet, row)

	index, _ := f.GetSheetIndex(sheet)
	f.SetActiveSheet(index)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write Excel file to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

// createComprehensiveStyles creates all Excel cell styles
func (g *ExcelGenerator) createComprehensiveStyles(f *excelize.File) *ExcelStyles {
	// Header style - cyan background, white text
	styleHeader, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#00B0F0"}},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Bold text style
	styleBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Basic border style
	styleBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Right-aligned border style
	styleBorderRight, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Center-aligned border style
	styleBorderCenter, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Number format - 2-4 decimals
	fourDecimalFormat := "#,##0.00##"

	// Currency format - Rp prefix, 2-4 decimals with left-aligned Rp and right-aligned amount
	currencyFormat := "\"Rp \"* #,##0.00##;-\"Rp \"* #,##0.00##;-"
	styleNumber, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Number format with borders
	styleNumberBorder, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Subtotal style - bold with borders
	styleSubtotal, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Right-aligned subtotal with numbers
	styleSubtotalRight, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Grand total style - bold, centered
	styleGrandTotal, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Right-aligned grand total with numbers
	styleGrandTotalRight, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Gray background for section titles
	styleGrayBackground, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#D9D9D9"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Header style without borders
	styleNoBorder, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Alternating row style - light gray
	styleAlternateRow, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#F2F2F2"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Center-aligned alternating row
	styleAlternateRowCenter, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#F2F2F2"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Number-formatted alternating row
	styleAlternateRowNumber, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#F2F2F2"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Right-aligned number format with borders
	styleNumberBorderRight, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &fourDecimalFormat,
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Right-aligned number format for alternating rows
	styleAlternateRowNumberRight, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &fourDecimalFormat,
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#F2F2F2"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Currency format with borders
	styleCurrencyBorder, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Currency format for alternating rows
	styleCurrencyAlternateRow, _ := f.NewStyle(&excelize.Style{
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#F2F2F2"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Bold currency format for headers
	styleCurrencyHeader, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		Alignment:    &excelize.Alignment{Horizontal: "left"},
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Normal currency format
	styleCurrencyNormal, _ := f.NewStyle(&excelize.Style{
		Alignment:    &excelize.Alignment{Horizontal: "left"},
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Bottom border styles
	styleBottomBorderOnly, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 0},
			{Type: "right", Color: "00B0F0", Style: 0},
			{Type: "top", Color: "00B0F0", Style: 0},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	styleBottomBorderBold, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 0},
			{Type: "right", Color: "00B0F0", Style: 0},
			{Type: "top", Color: "00B0F0", Style: 0},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	styleBottomBorderRight, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 0},
			{Type: "right", Color: "00B0F0", Style: 0},
			{Type: "top", Color: "00B0F0", Style: 0},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	styleBottomBorderNumber, _ := f.NewStyle(&excelize.Style{
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		CustomNumFmt: &fourDecimalFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 0},
			{Type: "right", Color: "00B0F0", Style: 0},
			{Type: "top", Color: "00B0F0", Style: 0},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	styleBottomBorderCurrency, _ := f.NewStyle(&excelize.Style{
		Alignment:    &excelize.Alignment{Horizontal: "right"},
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 0},
			{Type: "right", Color: "00B0F0", Style: 0},
			{Type: "top", Color: "00B0F0", Style: 0},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Bold currency format with borders
	styleCurrencyBorderBold, _ := f.NewStyle(&excelize.Style{
		Font:         &excelize.Font{Bold: true},
		CustomNumFmt: &currencyFormat,
		Fill:         excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Bold text with right alignment
	styleBoldRight, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "right"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	// Border styles for information section
	styleInfoBorderTop, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
		},
	})

	styleInfoBorderMiddle, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
		},
	})

	styleInfoBorderBottom, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Keep original for backward compatibility
	styleInfoBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "00B0F0", Style: 1},
			{Type: "right", Color: "00B0F0", Style: 1},
			{Type: "top", Color: "00B0F0", Style: 1},
			{Type: "bottom", Color: "00B0F0", Style: 1},
		},
	})

	// Normal text style
	styleNormal, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	return &ExcelStyles{
		styleHeader:                  styleHeader,
		styleBold:                    styleBold,
		styleBorder:                  styleBorder,
		styleBorderRight:             styleBorderRight,
		styleBorderCenter:            styleBorderCenter,
		styleNumber:                  styleNumber,
		styleNumberBorder:            styleNumberBorder,
		styleNumberBorderRight:       styleNumberBorderRight,
		styleSubtotal:                styleSubtotal,
		styleSubtotalRight:           styleSubtotalRight,
		styleGrandTotal:              styleGrandTotal,
		styleGrandTotalRight:         styleGrandTotalRight,
		styleGrayBackground:          styleGrayBackground,
		styleNoBorder:                styleNoBorder,
		styleAlternateRow:            styleAlternateRow,
		styleAlternateRowCenter:      styleAlternateRowCenter,
		styleAlternateRowNumber:      styleAlternateRowNumber,
		styleAlternateRowNumberRight: styleAlternateRowNumberRight,
		styleCurrencyBorder:          styleCurrencyBorder,
		styleCurrencyAlternateRow:    styleCurrencyAlternateRow,
		styleCurrencyHeader:          styleCurrencyHeader,
		styleCurrencyNormal:          styleCurrencyNormal,
		styleCurrencyBorderBold:      styleCurrencyBorderBold,
		styleBoldRight:               styleBoldRight,
		styleInfoBorder:              styleInfoBorder,
		styleInfoBorderTop:           styleInfoBorderTop,
		styleInfoBorderMiddle:        styleInfoBorderMiddle,
		styleInfoBorderBottom:        styleInfoBorderBottom,
		styleBottomBorderOnly:        styleBottomBorderOnly,
		styleBottomBorderBold:        styleBottomBorderBold,
		styleBottomBorderRight:       styleBottomBorderRight,
		styleBottomBorderNumber:      styleBottomBorderNumber,
		styleBottomBorderCurrency:    styleBottomBorderCurrency,
		styleNormal:                  styleNormal,
	}
}

// writeMainHeader writes the report header with logos and title
func (g *ExcelGenerator) writeMainHeader(f *excelize.File, sheet string, startRow int, data *crEntity.OJKCustomerReportDataDto, styles *ExcelStyles) int {
	row := startRow

	// Add logos and title
	// Add Sim logo
	logoPath := g.getLogoPath("logo.png")
	if err := f.AddPicture(sheet, fmt.Sprintf("C%d", row), logoPath, &excelize.GraphicOptions{
		ScaleX:          0.70,
		ScaleY:          0.70,
		OffsetX:         5,
		OffsetY:         5,
		LockAspectRatio: true,
		Positioning:     "oneCell",
	}); err != nil {
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Sim")
	}

	// Set report title
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), data.ReportTitle)
	f.MergeCell(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("G%d", row))
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("G%d", row), styles.styleNoBorder)

	ojkPath := g.getLogoPath("ojk.png")
	if err := f.AddPicture(sheet, fmt.Sprintf("I%d", row), ojkPath, &excelize.GraphicOptions{
		ScaleX:          0.70,
		ScaleY:          0.70,
		OffsetX:         5,
		OffsetY:         5,
		LockAspectRatio: true,
		Positioning:     "oneCell",
	}); err != nil {
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), "OJK")
	}
	row++

	// Add period subtitle
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("Periode : %s - %s",
		g.formatDate(data.Period.StartDate),
		g.formatDate(data.Period.EndDate)))

	periodStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	f.MergeCell(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("G%d", row))
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("G%d", row), periodStyle)
	row++

	return row
}

// writeCustomerInfoSection writes customer details
func (g *ExcelGenerator) writeCustomerInfoSection(f *excelize.File, sheet string, startRow int, data *crEntity.OJKCustomerReportDataDto, styles *ExcelStyles) int {
	row := startRow

	// Customer information
	var customerInfo []struct {
		label      string
		value      interface{}
		isCurrency bool
	}

	if data.Customer.Company != nil {
		// Customer with company - show Company ID
		customerInfo = []struct {
			label      string
			value      interface{}
			isCurrency bool
		}{
			{"ID Perusahaan", data.Customer.Company.CompanyCode, false},
			{"Nomor Kepesertaan", data.Customer.ID, false},
			{"Akumulasi Iuran", data.Summary.AccumulatedContribution, true},
			{"Akumulasi Hasil Pengembangan", data.Summary.AccumulatedDevelopmentResults, true},
			{"Akumulasi Biaya", -data.Summary.AccumulatedFees, true},
			{"Nominal Dana Kelolaan", data.Summary.ManagedFundValue, true},
		}
	} else {
		// Customer without company - exclude Company ID
		customerInfo = []struct {
			label      string
			value      interface{}
			isCurrency bool
		}{
			{"Nomor Kepesertaan", data.Customer.ID, false},
			{"Akumulasi Iuran", data.Summary.AccumulatedContribution, true},
			{"Akumulasi Hasil Pengembangan", data.Summary.AccumulatedDevelopmentResults, true},
			{"Akumulasi Biaya", -data.Summary.AccumulatedFees, true},
			{"Nominal Dana Kelolaan", data.Summary.ManagedFundValue, true},
		}
	}

	for _, item := range customerInfo {
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.label)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), ":")
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.value)

		// Apply styles based on data type
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("D%d", row), styles.styleNormal)
		if item.isCurrency {
			f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleCurrencyNormal)
		} else {
			f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleNormal)
		}
		row++
	}

	// Recipient information
	recipientRow := startRow
	f.SetCellValue(sheet, fmt.Sprintf("G%d", recipientRow), "Kepada Yth.")
	f.SetCellStyle(sheet, fmt.Sprintf("G%d", recipientRow), fmt.Sprintf("G%d", recipientRow), styles.styleBold)
	recipientRow++

	if data.Customer.Company != nil {
		recipientInfo := []string{
			fmt.Sprintf("Bapak/Ibu %s %s", data.Customer.FirstName, data.Customer.LastName),
			fmt.Sprintf("PT %s", data.Customer.Company.FirstName),
			data.Customer.Company.Address,
			fmt.Sprintf("Di %s", data.Customer.Company.Domisili),
		}

		for _, info := range recipientInfo {
			f.SetCellValue(sheet, fmt.Sprintf("G%d", recipientRow), info)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", recipientRow), fmt.Sprintf("G%d", recipientRow), styles.styleNormal)
			recipientRow++
		}
	}

	return row
}

// writeSideBySideSections writes IURAN, BIAYA, and PILIHAN INVESTASI sections
func (g *ExcelGenerator) writeSideBySideSections(f *excelize.File, sheet string, startRow int, data *crEntity.OJKCustomerReportDataDto, styles *ExcelStyles) int {
	row := startRow

	// IURAN section
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "• IURAN")
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleBold)
	row++

	// Add headers
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Jenis Iuran")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "Jumlah Iuran")
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("D%d", row), styles.styleHeader)
	row++

	// Add contribution data based on customer company status
	var iuranData []struct {
		label string
		value float64
	}
	var totalIuran float64 = data.TotalContribution

	if data.Customer.Company != nil {
		// Customer with company - show all contribution types
		iuranData = []struct {
			label string
			value float64
		}{
			{"1. Iuran Pemberi Kerja", data.ContributionSummary.EmployerContribution},
			{"2. Iuran Pekerja", data.ContributionSummary.EmployeeContribution},
			{"3. Iuran Sukarela Pekerja", data.ContributionSummary.VoluntaryContribution},
			{"4. Iuran Dana Pendidikan Anak", data.ContributionSummary.EducationFund},
		}
	} else {
		// Customer without company - only show Participant and Education Fund contributions
		iuranData = []struct {
			label string
			value float64
		}{
			{"1. Iuran Peserta", data.ContributionSummary.EmployeeContribution},
			{"2. Iuran Dana Pendidikan Anak", data.ContributionSummary.EducationFund},
		}
	}

	iuranStartRow := row
	for _, item := range iuranData {
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.label)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.value)
		f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleBorder)
		f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.styleCurrencyBorder)
		row++
	}

	// Add Total Iuran row
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Total Iuran")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), totalIuran)
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleBoldRight)
	f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.styleCurrencyBorderBold)
	row++

	// BIAYA section
	biayaRow := iuranStartRow - 2
	f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), "• BIAYA")
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("F%d", biayaRow), styles.styleBold)
	biayaRow++

	// Add headers
	f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), "Jenis Biaya")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", biayaRow), "Biaya")
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("G%d", biayaRow), styles.styleHeader)
	biayaRow++

	// Add fee data
	biayaData := []struct {
		label string
		value float64
	}{
		{"1. Administrasi Iuran", -data.FeeSummary.AdministrationFee},
		{"2. Pengelolaan Investasi", -data.FeeSummary.OperationalFee},
	}

	for _, item := range biayaData {
		f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), item.label)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", biayaRow), item.value)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("F%d", biayaRow), styles.styleBorder)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", biayaRow), fmt.Sprintf("G%d", biayaRow), styles.styleCurrencyBorder)
		biayaRow++
	}

	// PILIHAN INVESTASI section
	biayaRow += 2
	f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), "• PILIHAN INVESTASI")
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("F%d", biayaRow), styles.styleBold)
	biayaRow++

	// Add headers
	f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), "Jenis Dana")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", biayaRow), "(%)")
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("G%d", biayaRow), styles.styleHeader)
	biayaRow++

	// Add investment choice data
	for i, section := range data.TransactionSections {
		label := fmt.Sprintf("%d. %s", i+1, section.FundType)
		percentage := fmt.Sprintf("%d%%", section.Percentage)

		f.SetCellValue(sheet, fmt.Sprintf("F%d", biayaRow), label)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", biayaRow), percentage)
		f.SetCellStyle(sheet, fmt.Sprintf("F%d", biayaRow), fmt.Sprintf("F%d", biayaRow), styles.styleBorder)
		f.SetCellStyle(sheet, fmt.Sprintf("G%d", biayaRow), fmt.Sprintf("G%d", biayaRow), styles.styleBorderRight)
		biayaRow++
	}

	// Return the maximum row position to ensure content appears below both left and right sections
	maxRow := row
	if biayaRow > maxRow {
		maxRow = biayaRow
	}
	return maxRow
}

// writeFooterSection writes important information footer
func (g *ExcelGenerator) writeFooterSection(f *excelize.File, sheet string, startRow int, styles *ExcelStyles) int {
	row := startRow

	// Add important information
	importantInfo := []string{
		"Informasi Penting:",
		"1. Informasi dan data di dalam laporan transaksi ini bersifat rahasia.",
		"2. Laporan transaksi ini dibuat secara otomatis dan tidak membutuhkan tanda tangan resmi.",
		"3. Jika dalam 5 (lima) hari kerja setelah Laporan Transaksi ini diterima oleh Peserta dan DPLK SAM tidak menerima keberatan dari Peserta, maka Peserta dianggap setuju dengan",
		"    Laporan Transaksi ini. Jika ada perbedaan data antara laporan transaksi dan data DPLK SAM, maka data DPLK SAM yang dinyatakan benar. Untuk klarifikasi lebih lanjut",
		"    silahkan menghubungi DPLK SAM care pada alamat email dplk.sam@sinarmas-am.co.id",
		"4. Telpon : 021-50507000, Call center : 150555",
		"5. DPLK SAM terdaftar dan diawasi oleh Otoritas Jasa Keuangan.",
	}

	infoStartRow := row
	for _, info := range importantInfo {
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), info)
		row++
	}

	// Apply border to wrap entire information section and merge cells from C to I
	for i := infoStartRow; i < row; i++ {
		f.MergeCell(sheet, fmt.Sprintf("C%d", i), fmt.Sprintf("I%d", i))

		if i == infoStartRow {
			// First row (header) - top border with bold
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", i), fmt.Sprintf("I%d", i), styles.styleInfoBorderTop)
		} else if i == row-1 {
			// Last row - bottom border
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", i), fmt.Sprintf("I%d", i), styles.styleInfoBorderBottom)
		} else {
			// Middle rows - only left and right borders
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", i), fmt.Sprintf("I%d", i), styles.styleInfoBorderMiddle)
		}
	}

	return row
}

func (g *ExcelGenerator) writeTransactionSection(f *excelize.File, sheet string, startRow int, section *crEntity.TransactionSectionDto, styles *ExcelStyles) int {
	row := startRow

	// Add transaction headers
	headers := []string{"Jenis Transaksi", "Tanggal Transaksi", "Nilai Investasi (a)", "Tanggal Harga Unit", "Harga per Unit (b)", "Transaksi a/b (unit)", "Saldo (unit)"}
	for i, header := range headers {
		col := string(rune('C' + i)) // Start from column C
		f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), header)
		f.SetCellStyle(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("%s%d", col, row), styles.styleHeader)
	}
	row++

	// Add fund type title with gray background
	fundTypeTitle := fmt.Sprintf("Jenis Dana Investasi %s (%d%%)", section.FundType, section.Percentage)
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fundTypeTitle)
	f.MergeCell(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("I%d", row))
	f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("I%d", row), styles.styleGrayBackground)
	row++

	// Add transaction data with alternating colors
	for i, tx := range section.Transactions {
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), tx.TransactionType)

		// Handle opening balance differently
		if tx.TransactionType == "Saldo Awal" {
			// Show limited data for opening balance
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), g.formatDate(tx.TransactionDate))
			// Leave middle columns empty
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), "")                              // Nilai Investasi - empty
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "")                              // Tanggal Harga Unit - empty
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), "")                              // Harga per Unit - empty
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), "")                              // Transaksi a/b - empty
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), g.formatAmount(tx.BalanceUnits)) // Only show balance
		} else {
			// Show all fields for normal transactions
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), g.formatDate(tx.TransactionDate))
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), tx.InvestmentValue)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), g.formatDate(tx.UnitPriceDate))
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), tx.UnitPrice)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), g.formatAmount(tx.TransactionUnits))
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), g.formatAmount(tx.BalanceUnits))
		}

		// Apply alternating row colors
		if i%2 == 0 {
			// Even rows - white
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.styleBorderCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleCurrencyBorder)
			f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styles.styleBorderCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styles.styleCurrencyBorder)
			// Use right-aligned style for "Biaya Pengelolaan Investasi"
			if tx.TransactionType == "Biaya Pengelolaan Investasi" {
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styles.styleNumberBorderRight)
			} else {
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styles.styleNumberBorder)
			}
			f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styles.styleNumberBorder)
		} else {
			// Odd rows - gray
			f.SetCellStyle(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("C%d", row), styles.styleAlternateRow)
			f.SetCellStyle(sheet, fmt.Sprintf("D%d", row), fmt.Sprintf("D%d", row), styles.styleAlternateRowCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleCurrencyAlternateRow)
			f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("F%d", row), styles.styleAlternateRowCenter)
			f.SetCellStyle(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("G%d", row), styles.styleCurrencyAlternateRow)
			// Use right-aligned style for "Biaya Pengelolaan Investasi"
			if tx.TransactionType == "Biaya Pengelolaan Investasi" {
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styles.styleAlternateRowNumberRight)
			} else {
				f.SetCellStyle(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("H%d", row), styles.styleAlternateRowNumber)
			}
			f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), styles.styleAlternateRowNumber)
		}
		row++
	}

	// Add final balance summary
	// Set values for all columns
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), "Saldo akhir")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), section.FinalBalance)
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("I%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("J%d", row), " ")

	// Apply bottom borders
	for col := 'B'; col <= 'D'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderOnly)
	}
	// Format balance value
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleBottomBorderNumber)
	for col := 'F'; col <= 'J'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderBold)
	}
	row++

	// Add unit price row
	unitPriceDate := "31/05/2024" // Default fallback
	if len(section.Transactions) > 0 {
		unitPriceDate = g.formatDate(section.Transactions[len(section.Transactions)-1].UnitPriceDate)
	}
	// Set values for all columns
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("Harga unit pada %s", unitPriceDate))
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), section.UnitPrice)
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("I%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("J%d", row), " ")

	// Apply bottom borders
	for col := 'B'; col <= 'D'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderOnly)
	}
	// Format unit price value
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleBottomBorderCurrency)
	for col := 'F'; col <= 'J'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderOnly)
	}
	row++

	// Add investment value row
	// Set values for all columns
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("Nilai investasi Dana %s", strings.ToUpper(section.FundType)))
	// Merge cells for long text
	f.MergeCell(sheet, fmt.Sprintf("C%d", row), fmt.Sprintf("D%d", row))
	f.SetCellValue(sheet, fmt.Sprintf("E%d", row), section.FinalValue)
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("G%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("H%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("I%d", row), " ")
	f.SetCellValue(sheet, fmt.Sprintf("J%d", row), " ")

	// Apply bottom borders
	for col := 'B'; col <= 'D'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderOnly)
	}
	// Format investment value
	f.SetCellStyle(sheet, fmt.Sprintf("E%d", row), fmt.Sprintf("E%d", row), styles.styleBottomBorderCurrency)
	for col := 'F'; col <= 'J'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), styles.styleBottomBorderOnly)
	}
	row++

	return row
}

func (g *ExcelGenerator) formatDate(t time.Time) string {
	return t.Format("02/01/2006")
}

// formatAmount formats numeric values, returns "-" for zero
func (g *ExcelGenerator) formatAmount(amount float64) interface{} {
	if amount == 0 {
		return "-"
	}
	return amount
}

// getLogoPath returns absolute path to logo file
func (g *ExcelGenerator) getLogoPath(fileName string) string {
	pwd, err := os.Getwd()
	if err != nil {
		return filepath.Join("static", "images", fileName)
	}
	return filepath.Join(pwd, "static", "images", fileName)
}

// addReportBorderFrame creates complete border frame around report
func (g *ExcelGenerator) addReportBorderFrame(f *excelize.File, sheet string, endRow int) {
	// Calculate row positions
	topBorderRow := 2
	bottomBorderRow := endRow + 1
	finalPaddingRow := endRow + 2

	// Create border styles
	// Top-left corner
	topLeftCorner, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
		},
	})

	// Top-right corner
	topRightCorner, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "top", Color: "000000", Style: 2},
		},
	})

	// Bottom-left corner
	bottomLeftCorner, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	// Bottom-right corner
	bottomRightCorner, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	// Top border (middle sections)
	topBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 2},
		},
	})

	// Bottom border (middle sections)
	bottomBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	// Left border (middle sections)
	leftBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 2},
		},
	})

	// Right border (middle sections)
	rightBorder, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
		Border: []excelize.Border{
			{Type: "right", Color: "000000", Style: 2},
		},
	})

	// Apply top border
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", topBorderRow), fmt.Sprintf("B%d", topBorderRow), topLeftCorner)
	for col := 'C'; col <= 'I'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, topBorderRow), fmt.Sprintf("%c%d", col, topBorderRow), topBorder)
	}
	f.SetCellStyle(sheet, fmt.Sprintf("J%d", topBorderRow), fmt.Sprintf("J%d", topBorderRow), topRightCorner)

	// Apply side borders to content rows
	for row := topBorderRow + 1; row <= bottomBorderRow-1; row++ {
		f.SetCellStyle(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("B%d", row), leftBorder)
		f.SetCellStyle(sheet, fmt.Sprintf("J%d", row), fmt.Sprintf("J%d", row), rightBorder)
	}

	// Apply bottom border
	f.SetCellStyle(sheet, fmt.Sprintf("B%d", bottomBorderRow), fmt.Sprintf("B%d", bottomBorderRow), bottomLeftCorner)
	for col := 'C'; col <= 'I'; col++ {
		f.SetCellStyle(sheet, fmt.Sprintf("%c%d", col, bottomBorderRow), fmt.Sprintf("%c%d", col, bottomBorderRow), bottomBorder)
	}
	f.SetCellStyle(sheet, fmt.Sprintf("J%d", bottomBorderRow), fmt.Sprintf("J%d", bottomBorderRow), bottomRightCorner)

	// Fill border cells
	for col := 'C'; col <= 'I'; col++ {
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", col, topBorderRow), " ")
		f.SetCellValue(sheet, fmt.Sprintf("%c%d", col, bottomBorderRow), " ")
	}

	// Fill padding columns
	for row := topBorderRow; row <= finalPaddingRow; row++ {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), " ")
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), " ")
	}
}
