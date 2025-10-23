package handlers

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	crEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/company_report"
	customerReportEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/customer_report"
	reportEntity "github.com/raymondsugiarto/coffee-api/pkg/entity/report"
	companyparticipant "github.com/raymondsugiarto/coffee-api/pkg/module/report/company_participant"
	ojkcompanyreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_company_report"
	ojkcustomerreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_customer_report"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func GenerateOJKCompanyReportExcel(service ojkcompanyreport.OJKCompanyReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := new(crEntity.OJKCompanyReportFilterDto)
		if err := c.QueryParser(filter); err != nil {
			return status.New(status.BadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid query parameters"))
		}

		// Get company ID from context and set it in filter
		companyID := shared.GetCompanyID(c.Context())
		if companyID == nil {
			return status.New(status.Unauthorized, fiber.NewError(fiber.StatusUnauthorized, "company ID not found in context"))
		}
		filter.CompanyID = *companyID

		excelBytes, err := service.GenerateExcelOJKCompanyReportBytes(c.Context(), filter)
		if err != nil {
			return err
		}

		// Generate filename based on company and period
		monthStr := fmt.Sprintf("%04d-%02d", filter.Year, filter.Month)
		filename := "Company_Report_" + filter.CompanyID + "_" + monthStr + ".xlsx"

		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

		return c.SendStream(bytes.NewReader(excelBytes))
	}
}

func GenerateOJKCustomerReportExcel(service ojkcustomerreport.OJKCustomerReportService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := new(customerReportEntity.OJKCustomerReportFilterDto)
		if err := c.QueryParser(filter); err != nil {
			return status.New(status.BadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid query parameters"))
		}

		excelBytes, err := service.GenerateExcelForCompany(c.Context(), filter)
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
}

func GetCompanyParticipantReport(service companyparticipant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := new(reportEntity.CompanyParticipantFilter)
		if err := c.QueryParser(filter); err != nil {
			return status.New(status.BadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid query parameters"))
		}

		companyID := shared.GetCompanyID(c.Context())
		if companyID == nil {
			return status.New(status.Unauthorized, fiber.NewError(fiber.StatusUnauthorized, "company ID not found in context"))
		}
		filter.CompanyID = *companyID

		result, err := service.GetCompanyParticipantReport(c.Context(), filter)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
