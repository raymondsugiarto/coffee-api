package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/company"
	companyparticipant "github.com/raymondsugiarto/coffee-api/pkg/module/report/company_participant"
	ojkcompanyreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_company_report"
	ojkcustomerreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_customer_report"
)

func ReportRouter(app fiber.Router,
	ojkCompanyReportSvc ojkcompanyreport.OJKCompanyReportService,
	ojkCustomerReportSvc ojkcustomerreport.OJKCustomerReportService,
	companyParticipantSvc companyparticipant.Service) {
	// Transaction Report routes - filtered by company ID from context
	ojkGroup := app.Group("/transactions")
	ojkGroup.Get("/company", handlers.GenerateOJKCompanyReportExcel(ojkCompanyReportSvc))
	ojkGroup.Get("/customer", handlers.GenerateOJKCustomerReportExcel(ojkCustomerReportSvc))

	// Company Participant Report route
	app.Get("/participants", handlers.GetCompanyParticipantReport(companyParticipantSvc))
}
