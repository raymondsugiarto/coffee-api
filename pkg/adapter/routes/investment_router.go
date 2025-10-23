package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
)

func CompanyInvestmentRouter(app fiber.Router, svc investment.Service, itemService investmentitem.Service, paymentService investmentpayment.Service, distributionService investmentdistribution.Service) {
	app.Post("", handlers.CreateInvestment(svc))
	app.Get("", handlers.FindAllCompanyInvestment(svc))

	app.Get("/:id", handlers.FindInvestmentByID(svc))
	app.Delete("/:id", handlers.DeleteInvestment(svc))
	app.Put("/:id", handlers.UpdateInvestment(svc))

	app.Get("/distribution/company/:companyId", handlers.GetInvestmentDistributionByCompany(distributionService))
	app.Post("/distribution/create-update", handlers.CreateOrUpdateInvestmentDistribution(distributionService))

	app.Get("/monthly-contribution/company", handlers.GetTotalMonthlyCompanyContribution(paymentService))
	app.Get("/company/payments", handlers.FindAllInvestmentPayment(paymentService))
	app.Get("/payments/:id", handlers.FindInvestmentPaymentByID(paymentService))

	app.Get("/:investmentID/items", handlers.FindAllCompanyInvestmentItem(itemService))
	app.Get("/company/items", handlers.FindAllCompanyInvestmentItem(itemService))
	app.Get("/company/items/:id", handlers.FindInvestmentItemByID(itemService))
}

func CustomerInvestmentRouter(app fiber.Router, svc investment.Service, paymentService investmentpayment.Service) {
	app.Post("", handlers.CreateInvestment(svc))
	app.Get("", handlers.FindAllMyInvestment(svc))
	app.Post("/pay-contribution", handlers.CreateInvestmentFromPaymentFee(svc))

	app.Get("/:id", handlers.FindInvestmentByID(svc))
	app.Delete("/:id", handlers.DeleteInvestment(svc))
	app.Put("/:id", handlers.UpdateInvestment(svc))

	app.Get("/payments", handlers.FindAllMyInvestmentPayment(paymentService))
	app.Get("/payments/:id", handlers.FindInvestmentPaymentByID(paymentService))
	app.Post("/:id/payment", handlers.CreateInvestmentPayment(svc))
	app.Post("/:id/payment/:paymentId", handlers.UpdateInvestmentPaymentConfirmation(paymentService))
}
