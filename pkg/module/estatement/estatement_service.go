package estatement

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	b "github.com/getbrevo/brevo-go/lib"
	"github.com/gofiber/fiber/v2/log"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	ie "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
)

// Global mutex to prevent race conditions in the Maroto PDF library
var pdfGenerationMutex sync.Mutex

const (
	primaryColorRed   = 52
	primaryColorGreen = 73
	primaryColorBlue  = 94

	grayColorRed   = 240
	grayColorGreen = 240
	grayColorBlue  = 240

	lightGrayColorRed   = 250
	lightGrayColorGreen = 250
	lightGrayColorBlue  = 250

	footerGrayRed   = 128
	footerGrayGreen = 128
	footerGrayBlue  = 128

	whiteColorRed   = 255
	whiteColorGreen = 255
	whiteColorBlue  = 255
)

type EStatementService interface {
	GenerateEstatement(ctx context.Context, filter *entity.EstatementRequestDto) (*entity.EstatementDto, error)
	GeneratePDFEStatement(ctx context.Context, investmentItems []*ie.InvestmentStatementDto, startDate time.Time, endDate time.Time) (string, error)
	GenerateAndSendEmail(ctx context.Context, req *entity.EstatementEmailRequestDto) error
}

type estatementService struct {
	investmentItemSvc investmentitem.Service
	notificationSvc   notification.Service
	customerSvc       customer.Service
}

func NewEStatementService(investmentItemSvc investmentitem.Service, notificationSvc notification.Service, customerSvc customer.Service) EStatementService {
	return &estatementService{
		investmentItemSvc: investmentItemSvc,
		notificationSvc:   notificationSvc,
		customerSvc:       customerSvc,
	}
}

func (s *estatementService) validateDateRange(startDate, endDate time.Time) error {
	const maxMonths = 3

	if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	yearsDiff := endDate.Year() - startDate.Year()
	monthsDiff := int(endDate.Month()) - int(startDate.Month())
	totalMonthsDiff := yearsDiff*12 + monthsDiff

	if totalMonthsDiff == maxMonths && endDate.Day() > startDate.Day() {
		return fmt.Errorf("date range cannot exceed %d months", maxMonths)
	}

	if totalMonthsDiff > maxMonths {
		return fmt.Errorf("date range cannot exceed %d months", maxMonths)
	}

	return nil
}

func (s *estatementService) GenerateEstatement(ctx context.Context, filter *entity.EstatementRequestDto) (*entity.EstatementDto, error) {
	if err := s.validateDateRange(filter.StartDate, filter.EndDate); err != nil {
		return nil, err
	}

	items, err := s.investmentItemSvc.PrepareForStatement(ctx, &ie.InvestmentItemFindAllRequest{
		CustomerID:       filter.CustomerID,
		StartDate:        filter.StartDate,
		EndDate:          filter.EndDate,
		InvestmentStatus: model.InvestmentStatusSuccess,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to prepare investment items for statement: %w", err)
	}

	fileName, err := s.GeneratePDFEStatement(ctx, items, filter.StartDate, filter.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF e-statement: %w", err)
	}

	outputDir, err := s.ensureOutputDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get output directory: %w", err)
	}

	filePath := strings.Join([]string{outputDir, fileName}, "/")

	return &entity.EstatementDto{
		FilePath: filePath,
	}, nil
}

func (s *estatementService) ensureOutputDirectory() (string, error) {
	const (
		baseStorageDir   = "./storage"
		estatementSubDir = "estatement"
		dirPermission    = 0755
	)

	outputDir := strings.Join([]string{baseStorageDir, estatementSubDir}, "/")

	if err := os.MkdirAll(outputDir, dirPermission); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	return outputDir, nil
}

func (s *estatementService) GeneratePDFEStatement(ctx context.Context, investments []*ie.InvestmentStatementDto, startDate time.Time, endDate time.Time) (string, error) {
	const (
		pageNumberPattern = "Page {current} of {total}"
		pageMargin        = 10
		dirPermission     = 0755
		filePermission    = 0644
	)

	if len(investments) == 0 {
		return "", fmt.Errorf("no investment data provided")
	}

	firstInvestment := investments[0]
	fileName := s.generateFileName(firstInvestment.InvestmentItemDto)

	outputDir, err := s.ensureOutputDirectory()
	if err != nil {
		return "", fmt.Errorf("failed to ensure output directory: %w", err)
	}

	investmentsCopy := make([]*ie.InvestmentStatementDto, len(investments))
	copy(investmentsCopy, investments)

	go func(investmentItems []*ie.InvestmentStatementDto) {
		// Use global mutex to prevent race conditions in PDF library
		pdfGenerationMutex.Lock()
		defer pdfGenerationMutex.Unlock()

		cfg := config.NewBuilder().
			WithPageNumber(props.PageNumber{
				Pattern: pageNumberPattern,
				Place:   props.Bottom,
			}).
			WithLeftMargin(pageMargin).
			WithTopMargin(pageMargin).
			WithRightMargin(pageMargin).
			WithBottomMargin(pageMargin).
			WithPageSize(pagesize.A3).
			Build()

		m := maroto.New(cfg)

		s.addHeader(m, startDate, endDate)
		s.addCustomerInfo(m, firstInvestment.InvestmentItemDto)
		s.addInvestmentsTable(m, investmentItems)
		s.addNavInformation(m, investmentItems)
		s.addFeeInformation(m, investmentItems)
		s.addSummary(m, investmentItems)
		s.addFooter(m)

		pdf, err := m.Generate()
		if err != nil {
			fmt.Printf("failed to generate PDF %s: %v\n", fileName, err)
			return
		}

		pdfBytes := pdf.GetBytes()
		filePath := filepath.Join(outputDir, fileName)

		err = os.WriteFile(filePath, pdfBytes, filePermission)
		if err != nil {
			fmt.Printf("failed to save PDF %s: %v\n", fileName, err)
			return
		}

		fmt.Printf("PDF %s generated and saved successfully\n", fileName)
	}(investmentsCopy)

	return fileName, nil
}

func (s *estatementService) addHeader(m core.Maroto, startDate time.Time, endDate time.Time) {
	const (
		headerRowHeight     = 12
		periodRowHeight     = 8
		titleFontSize       = 16
		headerTitle         = "Investment E-Statement"
		topSpacing          = 3
		generatedLabel      = "Generated: %s"
		periodLabel         = "Periode: %s - %s"
		dateFormat          = "02/01/2006"
		periodDateFormat    = "02 January 2006"
		headerDateSpace     = 5
		fullColumnWidth     = 12
		separatorLineHeight = 5
		lineThickness       = 1
		sectionTopSpace     = 5
		sectionFontSize     = 14
		labelFontSize       = 10
		dataFontSize        = 10
		threeColumns        = 3
		fourColumns         = 4
		eightColumns        = 8
	)

	primaryColor := s.getPrimaryColor()

	m.AddRow(headerRowHeight,
		text.NewCol(eightColumns, headerTitle, props.Text{
			Size:  titleFontSize,
			Style: fontstyle.Bold,
			Align: align.Left,
			Color: primaryColor,
			Top:   topSpacing,
		}),
		text.NewCol(fourColumns, fmt.Sprintf(generatedLabel, time.Now().Format(dateFormat)), props.Text{
			Size:  labelFontSize,
			Align: align.Right,
			Top:   headerDateSpace,
		}),
	)

	// Add period information row
	periodText := fmt.Sprintf(periodLabel, startDate.Format(periodDateFormat), endDate.Format(periodDateFormat))
	m.AddRow(periodRowHeight,
		text.NewCol(fullColumnWidth, periodText, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Align: align.Center,
			Top:   topSpacing,
		}),
	)

	m.AddRow(separatorLineHeight, line.NewCol(fullColumnWidth, props.Line{
		Color:     primaryColor,
		Thickness: lineThickness,
	}))
}

func (s *estatementService) addCustomerInfo(m core.Maroto, investment *ie.InvestmentItemDto) {
	const (
		customerInfoTitle     = "Informasi Nasabah"
		customerInfoRowHeight = 8
		customerNameLabel     = "Nama:"
		customerIDLabel       = "ID Nasabah:"
		emailLabel            = "Email:"
		phoneLabel            = "Telepon:"
		notAvailableText      = "N/A"
		fullColumnWidth       = 12
		separatorLineHeight   = 5
		sectionTopSpace       = 5
		sectionFontSize       = 14
		labelFontSize         = 10
		dataFontSize          = 10
		threeColumns          = 3
	)

	m.AddRow(separatorLineHeight, col.New(fullColumnWidth))

	m.AddRow(sectionFontSize,
		text.NewCol(fullColumnWidth, customerInfoTitle, props.Text{
			Size:  sectionFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
	)

	customerName := notAvailableText
	customerID := notAvailableText
	customerEmail := notAvailableText
	customerPhone := notAvailableText

	if investment.Customer != nil {
		if investment.Customer.FirstName != "" {
			customerName = investment.Customer.FirstName
			if investment.Customer.LastName != "" {
				customerName += " " + investment.Customer.LastName
			}
		}
		if investment.Customer.Email != "" {
			customerEmail = investment.Customer.Email
		}
		if investment.Customer.PhoneNumber != "" {
			customerPhone = investment.Customer.PhoneNumber
		}
	}

	if investment.CustomerID != "" {
		customerID = investment.CustomerID
	}

	m.AddRow(customerInfoRowHeight,
		text.NewCol(threeColumns, customerNameLabel, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
		text.NewCol(threeColumns, customerName, props.Text{
			Size: dataFontSize,
			Top:  sectionTopSpace,
		}),
		text.NewCol(threeColumns, customerIDLabel, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
		text.NewCol(threeColumns, customerID, props.Text{
			Size: dataFontSize,
			Top:  sectionTopSpace,
		}),
	)

	m.AddRow(customerInfoRowHeight,
		text.NewCol(threeColumns, emailLabel, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
		text.NewCol(threeColumns, customerEmail, props.Text{
			Size: dataFontSize,
			Top:  sectionTopSpace,
		}),
		text.NewCol(threeColumns, phoneLabel, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
		text.NewCol(threeColumns, customerPhone, props.Text{
			Size: dataFontSize,
			Top:  sectionTopSpace,
		}),
	)
}

func (s *estatementService) addInvestmentsTable(m core.Maroto, items []*ie.InvestmentStatementDto) {

	const (
		investmentDetailsTitle  = "Detail Investasi"
		sectionTitleHeight      = 15
		tableRowHeight          = 12
		dateColumnHeader        = "Tanggal"
		codeColumnHeader        = "Kode"
		productNameColumnHeader = "Nama Produk"
		typeColumnHeader        = "Tipe"
		modalInvestasiHeader    = "Total Iuran"
		unitIPHeader            = "Total Unit"
		navTerkiniHeader        = "NAV Terkini"
		nilaiSaatIniHeader      = "Dana Peserta"
		tableHeaderFont         = 11
		tableCellSpace          = 3
		tableDataFont           = 10
		tableRightPadding       = 8
		tableRightPaddingSmall  = 4
		dateFormat              = "02/01/2006"
		notAvailable            = "N/A"
		evenRowIndex            = 0
		fullColumnWidth         = 12
		separatorLineHeight     = 5
		sectionTopSpace         = 5
		sectionFontSize         = 14
		oneColumn               = 1
		twoColumns              = 2
		threeColumns            = 3
	)

	headerColor := s.getPrimaryColor()
	grayColor := s.getGrayColor()

	m.AddRow(separatorLineHeight, col.New(fullColumnWidth))

	m.AddRow(sectionTitleHeight,
		text.NewCol(fullColumnWidth, investmentDetailsTitle, props.Text{
			Size:  sectionFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
	)

	headerRow := m.AddRow(tableRowHeight,
		text.NewCol(oneColumn, dateColumnHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
		}),
		text.NewCol(oneColumn, codeColumnHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
		}),
		text.NewCol(twoColumns, productNameColumnHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
		}),
		text.NewCol(oneColumn, typeColumnHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Center,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
		}),
		text.NewCol(twoColumns, modalInvestasiHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Right,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
			Right: tableRightPadding,
		}),
		text.NewCol(oneColumn, unitIPHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Right,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
			Right: tableRightPadding,
		}),
		text.NewCol(twoColumns, navTerkiniHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Right,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
			Right: tableRightPadding,
		}),
		text.NewCol(twoColumns, nilaiSaatIniHeader, props.Text{
			Size:  tableHeaderFont,
			Style: fontstyle.Bold,
			Align: align.Right,
			Color: &props.WhiteColor,
			Top:   tableCellSpace,
			Right: tableRightPaddingSmall,
		}),
	)
	headerRow.WithStyle(&props.Cell{
		BackgroundColor: headerColor,
	})

	for i, item := range items {
		investmentDate := item.InvestmentAt.Format(dateFormat)
		investmentCode := ""
		if item.Investment != nil {
			investmentCode = item.Investment.Code
		}
		if investmentCode == "" {
			investmentCode = notAvailable
		}

		productName := notAvailable
		if item.InvestmentProduct != nil && item.InvestmentProduct.Name != "" {
			productName = item.InvestmentProduct.Name
		}

		investmentType := string(item.Type)
		if investmentType == "" {
			investmentType = notAvailable
		}

		modalInvestasi := utils.FormatIndonesianNumber(item.TotalAmount)
		unit := fmt.Sprintf("%.4f", item.Unit)
		navTerkini := fmt.Sprintf("%.4f", item.CurrentNavAmount)
		nilaiSaatIni := utils.FormatIndonesianNumber(item.CurrentValue)

		rowColor := s.getWhiteColor()
		if i%2 == evenRowIndex {
			rowColor = grayColor
		}

		dataRow := m.AddRow(tableRowHeight,
			text.NewCol(oneColumn, investmentDate, props.Text{
				Size:  tableDataFont,
				Align: align.Center,
				Top:   tableCellSpace,
			}),
			text.NewCol(oneColumn, investmentCode, props.Text{
				Size:  tableDataFont,
				Align: align.Center,
				Top:   tableCellSpace,
			}),
			text.NewCol(twoColumns, productName, props.Text{
				Size:  tableDataFont,
				Align: align.Center,
				Top:   tableCellSpace,
			}),
			text.NewCol(oneColumn, investmentType, props.Text{
				Size:  tableDataFont,
				Align: align.Center,
				Top:   tableCellSpace,
			}),
			text.NewCol(twoColumns, modalInvestasi, props.Text{
				Size:  tableDataFont,
				Align: align.Right,
				Top:   tableCellSpace,
				Right: tableRightPadding,
			}),
			text.NewCol(oneColumn, unit, props.Text{
				Size:  tableDataFont,
				Align: align.Right,
				Top:   tableCellSpace,
				Right: tableRightPadding,
			}),
			text.NewCol(twoColumns, navTerkini, props.Text{
				Size:  tableDataFont,
				Align: align.Right,
				Top:   tableCellSpace,
				Right: tableRightPadding,
			}),
			text.NewCol(twoColumns, nilaiSaatIni, props.Text{
				Size:  tableDataFont,
				Align: align.Right,
				Top:   tableCellSpace,
				Right: tableRightPaddingSmall,
			}),
		)
		dataRow.WithStyle(&props.Cell{
			BackgroundColor: rowColor,
		})
	}
}

type MonthlyFee struct {
	Month            string
	YearMonth        string
	TotalAmount      float64
	TransactionCount int
}

func (s *estatementService) addFeeInformation(m core.Maroto, investments []*ie.InvestmentStatementDto) {
	const (
		feeInformationTitle = "Informasi Biaya"
		sectionTitleHeight  = 15
		feeRowHeight        = 12
		totalFeeLabel       = "Total Biaya Admin:"
		transactionFeeLabel = "Rincian Biaya Transaksi:"
		periodHeader        = "Periode"
		transactionHeader   = "Jumlah Transaksi"
		amountHeader        = "Total Biaya"
		tableCellSpace      = 3
		fullColumnWidth     = 12
		separatorLineHeight = 5
		sectionTopSpace     = 5
		sectionFontSize     = 14
		labelFontSize       = 12
		dataFontSize        = 10
		sixColumns          = 6
		fourColumns         = 4
	)

	var totalFeeAmount float64
	for _, inv := range investments {
		totalFeeAmount += inv.FeeAmount
	}

	grayColor := s.getGrayColor()
	lightGrayColor := s.getLightGrayColor()

	m.AddRow(separatorLineHeight, col.New(fullColumnWidth))

	m.AddRow(sectionTitleHeight,
		text.NewCol(fullColumnWidth, feeInformationTitle, props.Text{
			Size:  sectionFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
	)

	totalFeeRow := m.AddRow(feeRowHeight,
		text.NewCol(sixColumns, totalFeeLabel, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Style: fontstyle.Bold,
			Align: align.Left,
			Left:  sectionTopSpace,
		}),
		text.NewCol(sixColumns, utils.FormatIndonesianNumber(totalFeeAmount), props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Align: align.Right,
			Right: sectionTopSpace,
		}),
	)
	totalFeeRow.WithStyle(&props.Cell{
		BackgroundColor: grayColor,
	})

	m.AddRow(feeRowHeight,
		text.NewCol(fullColumnWidth, transactionFeeLabel, props.Text{
			Size:  labelFontSize,
			Style: fontstyle.Bold,
			Top:   tableCellSpace,
			Left:  sectionTopSpace,
		}),
	)

	// Group investments by month
	monthlyFees := make(map[string]*MonthlyFee)
	for _, inv := range investments {
		if inv.FeeAmount > 0 {
			yearMonth := inv.InvestmentAt.Format("2006-01")
			monthName := inv.InvestmentAt.Format("January 2006")

			if _, exists := monthlyFees[yearMonth]; !exists {
				monthlyFees[yearMonth] = &MonthlyFee{
					Month:            monthName,
					YearMonth:        yearMonth,
					TotalAmount:      0,
					TransactionCount: 0,
				}
			}

			monthlyFees[yearMonth].TotalAmount += inv.FeeAmount
			monthlyFees[yearMonth].TransactionCount++
		}
	}

	// Convert map to slice for sorting
	var sortedFees []*MonthlyFee
	for _, fee := range monthlyFees {
		sortedFees = append(sortedFees, fee)
	}

	// Sort by year-month
	for i := 0; i < len(sortedFees); i++ {
		for j := i + 1; j < len(sortedFees); j++ {
			if sortedFees[i].YearMonth > sortedFees[j].YearMonth {
				sortedFees[i], sortedFees[j] = sortedFees[j], sortedFees[i]
			}
		}
	}

	// Display header if there are monthly fees
	if len(sortedFees) > 0 {
		headerRow := m.AddRow(feeRowHeight,
			text.NewCol(fourColumns, periodHeader, props.Text{
				Size:  labelFontSize,
				Style: fontstyle.Bold,
				Align: align.Center,
				Color: &props.WhiteColor,
				Top:   tableCellSpace,
			}),
			text.NewCol(fourColumns, transactionHeader, props.Text{
				Size:  labelFontSize,
				Style: fontstyle.Bold,
				Align: align.Center,
				Color: &props.WhiteColor,
				Top:   tableCellSpace,
			}),
			text.NewCol(fourColumns, amountHeader, props.Text{
				Size:  labelFontSize,
				Style: fontstyle.Bold,
				Align: align.Center,
				Color: &props.WhiteColor,
				Top:   tableCellSpace,
			}),
		)
		headerRow.WithStyle(&props.Cell{
			BackgroundColor: s.getPrimaryColor(),
		})

		// Display monthly fee data
		for i, fee := range sortedFees {
			rowColor := s.getWhiteColor()
			if i%2 == 0 {
				rowColor = lightGrayColor
			}

			feeRow := m.AddRow(feeRowHeight,
				text.NewCol(fourColumns, fee.Month, props.Text{
					Size:  dataFontSize,
					Top:   tableCellSpace,
					Align: align.Center,
				}),
				text.NewCol(fourColumns, fmt.Sprintf("%d transaksi", fee.TransactionCount), props.Text{
					Size:  dataFontSize,
					Top:   tableCellSpace,
					Align: align.Center,
				}),
				text.NewCol(fourColumns, utils.FormatIndonesianNumber(fee.TotalAmount), props.Text{
					Size:  dataFontSize,
					Top:   tableCellSpace,
					Align: align.Right,
					Right: sectionTopSpace,
				}),
			)
			feeRow.WithStyle(&props.Cell{
				BackgroundColor: rowColor,
			})
		}
	}
}

func (s *estatementService) addSummary(m core.Maroto, investments []*ie.InvestmentStatementDto) {

	const (
		investmentSummaryTitle = "Ringkasan Investasi"
		sectionTitleHeight     = 15
		summaryRowHeight       = 12
		totalInvestmentsLabel  = "Total Investasi:"
		modalInvestasiLabel    = "Total Iuran:"
		nilaiSaatIniLabel      = "Total Dana Peserta:"
		gainLossLabel          = "Total Hasil Investasi:"
		labelFontSize          = 12
		tableCellSpace         = 3
		fullColumnWidth        = 12
		separatorLineHeight    = 5
		sectionTopSpace        = 6
		sectionFontSize        = 14
		sixColumns             = 6
	)

	var totalModalInvestasi, totalNilaiSaatIni, totalGainLoss float64
	totalCount := len(investments)

	for _, inv := range investments {
		totalModalInvestasi += inv.TotalAmount
		totalNilaiSaatIni += inv.CurrentValue
		totalGainLoss += inv.GainLoss
	}

	var gainLossPercentage float64
	if totalModalInvestasi > 0 {
		gainLossPercentage = (totalGainLoss / totalModalInvestasi) * 100
	}

	grayColor := s.getGrayColor()

	m.AddRow(separatorLineHeight, col.New(fullColumnWidth))

	m.AddRow(sectionTitleHeight,
		text.NewCol(fullColumnWidth, investmentSummaryTitle, props.Text{
			Size:  sectionFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
	)

	totalCountRow := m.AddRow(summaryRowHeight,
		text.NewCol(sixColumns, totalInvestmentsLabel, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Style: fontstyle.Bold,
			Align: align.Left,
			Left:  sectionTopSpace,
		}),
		text.NewCol(sixColumns, fmt.Sprintf("%d", totalCount), props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Align: align.Right,
			Right: sectionTopSpace,
		}),
	)

	totalCountRow.WithStyle(&props.Cell{
		BackgroundColor: grayColor,
	})

	modalInvestasiRow := m.AddRow(summaryRowHeight,
		text.NewCol(sixColumns, modalInvestasiLabel, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Style: fontstyle.Bold,
			Align: align.Left,
			Left:  sectionTopSpace,
		}),
		text.NewCol(sixColumns, utils.FormatIndonesianNumber(totalModalInvestasi), props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Align: align.Right,
			Right: sectionTopSpace,
		}),
	)
	modalInvestasiRow.WithStyle(&props.Cell{
		BackgroundColor: s.getLightGrayColor(),
	})

	nilaiSaatIniRow := m.AddRow(summaryRowHeight,
		text.NewCol(sixColumns, nilaiSaatIniLabel, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Style: fontstyle.Bold,
			Align: align.Left,
			Left:  sectionTopSpace,
		}),
		text.NewCol(sixColumns, utils.FormatIndonesianNumber(totalNilaiSaatIni), props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Align: align.Right,
			Right: sectionTopSpace,
		}),
	)
	nilaiSaatIniRow.WithStyle(&props.Cell{
		BackgroundColor: s.getLightGrayColor(),
	})

	gainLossText := fmt.Sprintf("%s (%.1f%%)", utils.FormatIndonesianNumber(totalGainLoss), gainLossPercentage)
	gainLossRow := m.AddRow(summaryRowHeight,
		text.NewCol(sixColumns, gainLossLabel, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Style: fontstyle.Bold,
			Align: align.Left,
			Left:  sectionTopSpace,
		}),
		text.NewCol(sixColumns, gainLossText, props.Text{
			Size:  labelFontSize,
			Top:   tableCellSpace,
			Align: align.Right,
			Right: sectionTopSpace,
		}),
	)
	gainLossRow.WithStyle(&props.Cell{
		BackgroundColor: s.getLightGrayColor(),
	})
}

func (s *estatementService) addNavInformation(m core.Maroto, investments []*ie.InvestmentStatementDto) {
	const (
		navInformationTitle = "Informasi NAV"
		sectionTitleHeight  = 15
		navRowHeight        = 10
		labelFontSize       = 12
		dataFontSize        = 10
		fullColumnWidth     = 12
		separatorLineHeight = 5
		sectionTopSpace     = 5
		sectionFontSize     = 14
		bulletPoint         = "• %s: %s per unit (%s)"
	)

	// Get unique products with their latest NAV and dates
	type NavInfo struct {
		Name   string
		Amount float64
		Date   time.Time
	}

	productNavMap := make(map[string]NavInfo)

	for _, inv := range investments {
		if inv.CurrentNavAmount > 0 && inv.InvestmentProduct != nil {
			productNavMap[inv.InvestmentProductID] = NavInfo{
				Name:   inv.InvestmentProduct.Name,
				Amount: inv.CurrentNavAmount,
				Date:   inv.CurrentNavDate,
			}
		}
	}

	if len(productNavMap) == 0 {
		return // Skip if no NAV data
	}

	m.AddRow(separatorLineHeight, col.New(fullColumnWidth))

	m.AddRow(sectionTitleHeight,
		text.NewCol(fullColumnWidth, navInformationTitle, props.Text{
			Size:  sectionFontSize,
			Style: fontstyle.Bold,
			Top:   sectionTopSpace,
		}),
	)

	for _, navInfo := range productNavMap {
		productName := navInfo.Name
		if productName == "" {
			productName = "N/A"
		}

		navText := fmt.Sprintf(bulletPoint,
			productName,
			fmt.Sprintf("%.4f", navInfo.Amount),
			navInfo.Date.Format("02/01/2006"))

		m.AddRow(navRowHeight,
			text.NewCol(fullColumnWidth, navText, props.Text{
				Size: dataFontSize,
				Top:  sectionTopSpace,
				Left: sectionTopSpace,
			}),
		)
	}
}

func (s *estatementService) addFooter(m core.Maroto) {
	const (
		footerSpacing    = 10
		footerLineHeight = 3
		footerTextHeight = 12
		companyName      = "DPLK Sinarmas Asset Managament"
		generatedOnLabel = "Dibuat pada: %s"
		disclaimerText   = "Dokumen ini dibuat secara komputerisasi dan tidak memerlukan tanda tangan."
		contactInfo      = "Hubungi Kami: info@dplk.com | +62-21-1234567"
		footerLabelFont  = 12
		footerSmallFont  = 8
		footerTopSpace   = 2
		dateTimeFormat   = "02/01/2006 15:04:05"
		lineThickness    = 1
		tableCellSpace   = 3
		fullColumnWidth  = 12
		sixColumns       = 6
	)

	primaryColor := s.getPrimaryColor()
	lightGray := s.getFooterGrayColor()

	m.AddRow(footerSpacing, col.New(fullColumnWidth))

	m.AddRow(footerLineHeight, line.NewCol(fullColumnWidth, props.Line{
		Color:     primaryColor,
		Thickness: lineThickness,
	}))

	m.AddRow(footerTextHeight,
		text.NewCol(sixColumns, companyName, props.Text{
			Size:  footerLabelFont,
			Style: fontstyle.Bold,
			Color: primaryColor,
			Top:   footerTopSpace,
		}),
		text.NewCol(sixColumns, fmt.Sprintf(generatedOnLabel, time.Now().Format(dateTimeFormat)), props.Text{
			Size:  footerSmallFont,
			Align: align.Right,
			Color: lightGray,
			Top:   footerTopSpace,
		}),
	)

	m.AddRow(footerSpacing,
		text.NewCol(fullColumnWidth, disclaimerText, props.Text{
			Size:  footerSmallFont,
			Align: align.Center,
			Color: lightGray,
			Top:   footerTopSpace,
		}),
	)

	m.AddRow(footerTextHeight,
		text.NewCol(fullColumnWidth, contactInfo, props.Text{
			Size:  footerSmallFont,
			Align: align.Center,
			Color: lightGray,
			Top:   tableCellSpace,
		}),
	)
}

func (s *estatementService) generateFileName(investment *ie.InvestmentItemDto) string {
	const (
		unknownCustomer    = "unknown"
		fileNameDateFormat = "20060102"
		fileNameFormat     = "EStatement_%s_%s.pdf"
	)

	var customerIdentifier string

	if investment.Customer != nil {
		if investment.Customer.FirstName != "" {
			customerIdentifier = investment.Customer.FirstName
			if investment.Customer.LastName != "" {
				customerIdentifier += "_" + investment.Customer.LastName
			}
		} else {
			customerIdentifier = investment.CustomerID
		}

		customerIdentifier = strings.ReplaceAll(customerIdentifier, " ", "_")
	} else if investment.CustomerID != "" {
		customerIdentifier = investment.CustomerID
	} else {
		customerIdentifier = unknownCustomer
	}

	currentDate := time.Now().Format(fileNameDateFormat)
	return fmt.Sprintf(fileNameFormat, customerIdentifier, currentDate)
}

func (s *estatementService) getPrimaryColor() *props.Color {
	return &props.Color{Red: primaryColorRed, Green: primaryColorGreen, Blue: primaryColorBlue}
}

func (s *estatementService) getGrayColor() *props.Color {
	return &props.Color{Red: grayColorRed, Green: grayColorGreen, Blue: grayColorBlue}
}

func (s *estatementService) getLightGrayColor() *props.Color {
	return &props.Color{Red: lightGrayColorRed, Green: lightGrayColorGreen, Blue: lightGrayColorBlue}
}

func (s *estatementService) getWhiteColor() *props.Color {
	return &props.Color{Red: whiteColorRed, Green: whiteColorGreen, Blue: whiteColorBlue}
}

func (s *estatementService) getFooterGrayColor() *props.Color {
	return &props.Color{Red: footerGrayRed, Green: footerGrayGreen, Blue: footerGrayBlue}
}

func (s *estatementService) GenerateAndSendEmail(ctx context.Context, req *entity.EstatementEmailRequestDto) error {
	log.Infof("Generating and sending e-statement email for customer %s", req.CustomerID)

	if err := s.validateDateRange(req.StartDate, req.EndDate); err != nil {
		log.Errorf("Date range validation failed: %v", err)
		return err
	}

	customerData, err := s.customerSvc.FindByIDWithScope(ctx, req.CustomerID, []string{"complete"})
	if err != nil {
		log.Errorf("Failed to find customer %s: %v", req.CustomerID, err)
		return fmt.Errorf("failed to find customer: %w", err)
	}

	if customerData.Email == "" {
		log.Errorf("Customer %s has no email address", req.CustomerID)
		return fmt.Errorf("customer has no email address")
	}

	items, err := s.investmentItemSvc.PrepareForStatement(ctx, &ie.InvestmentItemFindAllRequest{
		CustomerID:       req.CustomerID,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		InvestmentStatus: model.InvestmentStatusSuccess,
	})

	if err != nil {
		log.Errorf("Failed to prepare investment items for statement: %v", err)
		return fmt.Errorf("failed to prepare investment items for statement: %w", err)
	}

	fileName, err := s.GeneratePDFEStatement(ctx, items, req.StartDate, req.EndDate)
	if err != nil {
		log.Errorf("Failed to generate PDF e-statement: %v", err)
		return fmt.Errorf("failed to generate PDF e-statement: %w", err)
	}

	outputDir, err := s.ensureOutputDirectory()
	if err != nil {
		log.Errorf("Failed to get output directory: %v", err)
		return fmt.Errorf("failed to get output directory: %w", err)
	}

	filePath := strings.Join([]string{outputDir, fileName}, "/")

	customerName := customerData.FirstName
	if customerData.LastName != "" {
		customerName += " " + customerData.LastName
	}
	if customerName == "" {
		customerName = "Customer"
	}

	emailReq := &entity.NotificationInputDto{
		To: []b.SendSmtpEmailTo{
			{
				Email: customerData.Email,
				Name:  customerName,
			},
		},
		TemplateID: 1,
		Data: map[string]interface{}{
			"customer_name": customerName,
			"start_date":    req.StartDate.Format("02/01/2006"),
			"end_date":      req.EndDate.Format("02/01/2006"),
			"file_path":     filePath,
		},
	}

	_, err = s.notificationSvc.SendEmailTemplate(ctx, emailReq)
	if err != nil {
		log.Errorf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Infof("E-statement email sent successfully to %s", customerData.Email)
	return nil
}
