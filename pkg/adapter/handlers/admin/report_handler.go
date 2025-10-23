package handlers

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/company_report"
	customerReportEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer_report"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/aum"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/channel"
	contributionsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/contribution_summary"
	ojkcompanyreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_company_report"
	ojkcustomerreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_customer_report"
	participantsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/participant_summary"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/portfolio"
	summaryaum "github.com/raymondsugiarto/coffee-api/pkg/module/report/summary_aum"
	transactionhistory "github.com/raymondsugiarto/coffee-api/pkg/module/report/transaction_history"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func ReportChannel(service channel.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ReportTransactionChannelFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateReportChannel(c.Context(), req)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", "attachment; filename=riwayat_transaksi.xlsx")
			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetTransactionReportChannel(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func ReportParticipantSummary(service participantsummary.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ReportParticipantSummaryFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		// Jika Accept adalah Excel, kirim Excel
		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateReportParticipantSummary(c.Context(), req)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="download.xlsx"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		// Jika bukan Excel, kirim JSON
		result, err := service.GetParticipantSummary(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func ReportContributionSummary(service contributionsummary.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ReportContributionSummaryFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateContributionReport(c.Context(), req)
			if err != nil {
				return err
			}
			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="download.xlsx"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetContributionSummary(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func ReportSummaryAum(service summaryaum.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ReportSummaryAumFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateReportSummaryAum(c.Context(), req)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="download.xlsx"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetSummaryAum(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func formatIndonesianMonthForFilename(month int) string {
	months := [...]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return "Unknown"
}

func ReportAUM(service aum.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ReportAUMFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateReportAUM(c.Context(), req)
			if err != nil {
				return err
			}

			monthName := formatIndonesianMonthForFilename(req.Month)
			filename := fmt.Sprintf("Report_AUM_%s_%s_%d.xlsx", string(req.CompanyType), monthName, req.Year)
			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetReportAum(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GenerateOJKCompanyReportExcel(service ojkcompanyreport.OJKCompanyReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := new(crEntity.OJKCompanyReportFilterDto)
		if err := c.QueryParser(filter); err != nil {
			return status.New(status.BadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid query parameters"))
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {

			excelBytes, err := service.GenerateExcelOJKCompanyReportBytes(c.Context(), filter)
			if err != nil {
				return err
			}

			// Generate filename based on company and period
			monthStr := fmt.Sprintf("%04d-%02d", filter.Year, filter.Month)
			filename := "Company_Report_" + monthStr + ".xlsx"
			if filter.CompanyID != "" {
				filename = "Company_Report_" + filter.CompanyID + "_" + monthStr + ".xlsx"
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetOJKCompanyReportData(c.Context(), filter)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func ReportTransactionHistory(service transactionhistory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.TransactionHistoryFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.GetTransactionHistoryReport(c.Context(), req)
		if err != nil {
			return err
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateExcel(c.Context(), req, result)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", "attachment; filename=riwayat_transaksi.xlsx")
			return c.SendStream(bytes.NewReader(excelBytes))
		}

		return c.JSON(result)
	}
}

func GenerateOJKCustomerReportExcel(service ojkcustomerreport.OJKCustomerReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := new(customerReportEntity.OJKCustomerReportFilterDto)
		if err := c.QueryParser(filter); err != nil {
			return status.New(status.BadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid query parameters"))
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {

			excelBytes, err := service.GenerateExcelOJKCustomerReportBytes(c.Context(), filter)
			if err != nil {
				return err
			}

			// Generate filename based on customer and period
			monthStr := fmt.Sprintf("%04d-%02d", filter.Year, filter.Month)
			filename := "Customer_Report_" + monthStr + ".xlsx"
			if filter.CustomerID != "" {
				filename = "Customer_Report_" + filter.CustomerID + "_" + monthStr + ".xlsx"
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

			return c.SendStream(bytes.NewReader(excelBytes))
		}

		result, err := service.GetOJKCustomerReportData(c.Context(), filter)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func ReportPortfolio(service portfolio.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.PortfolioFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAllPortfolioWithNav(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
