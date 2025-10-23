package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/aum"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/channel"
	contributionsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/contribution_summary"
	ojkcompanyreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_company_report"
	ojkcustomerreport "github.com/raymondsugiarto/coffee-api/pkg/module/report/ojk_customer_report"
	participantsummary "github.com/raymondsugiarto/coffee-api/pkg/module/report/participant_summary"
	"github.com/raymondsugiarto/coffee-api/pkg/module/report/portfolio"
	summaryaum "github.com/raymondsugiarto/coffee-api/pkg/module/report/summary_aum"
	transactionhistory "github.com/raymondsugiarto/coffee-api/pkg/module/report/transaction_history"
)

func ReportRouter(app fiber.Router, channelSvc channel.Service, participantSvc participantsummary.Service, contributionSvc contributionsummary.Service, summaryaumSvc summaryaum.Service, aumSvc aum.Service, ojkCompanyReportSvc ojkcompanyreport.OJKCompanyReportService, ojkCustomerReportSvc ojkcustomerreport.OJKCustomerReportService, transactionHistorySvc transactionhistory.Service, portfolioSvc portfolio.Service) {
	app.Get("/channels", handlers.ReportChannel(channelSvc))
	app.Get("/participant-summary", handlers.ReportParticipantSummary(participantSvc))
	app.Get("/contribution-summary", handlers.ReportContributionSummary(contributionSvc))
	app.Get("/summary-aum", handlers.ReportSummaryAum(summaryaumSvc))
	app.Get("/aum", handlers.ReportAUM(aumSvc))
	app.Get("/transaction-history", handlers.ReportTransactionHistory(transactionHistorySvc))
	app.Get("/portfolios", handlers.ReportPortfolio(portfolioSvc))

	// Transaction Report routes
	ojkGroup := app.Group("/transactions")
	ojkGroup.Get("/company", handlers.GenerateOJKCompanyReportExcel(ojkCompanyReportSvc))
	ojkGroup.Get("/customer", handlers.GenerateOJKCustomerReportExcel(ojkCustomerReportSvc))
}
