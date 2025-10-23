package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
)

func AdminDashboardRouter(app fiber.Router,
	customerService customer.Service,
	investmentProductService investmentproduct.Service,
	unitLinkService unitlink.Service,
	companyService company.Service,
) {
	app.Get("/customers/summary", handlers.GetCountCustomerByType(customerService))
	app.Get("/companies/summary", handlers.GetCountCompanyByType(companyService))
	app.Get("/investment-distributions/summary", handlers.InvestmentProductSummary(investmentProductService))
	app.Get("/unit-links/summary", handlers.UnitLinkSummaryPerType(unitLinkService))
}

func CompanyDashboardRouter(app fiber.Router, customerService customer.Service, unitLinkService unitlink.Service, investmentDistributionService investmentdistribution.Service, investmentPaymentService investmentpayment.Service) {
	app.Get("/customers/summary", handlers.GetCountCustomerByType(customerService))
	app.Get("/new-customers/summary", handlers.GetCountNewCustomerByType(customerService))
	app.Get("/unit-links/summary", handlers.UnitLinkSummaryByCompany(unitLinkService))
	app.Get("/investment-distributions/summary", handlers.InvestmentDistributionSummaryByCompany(investmentDistributionService))
	app.Get("/investment-payments/summary", handlers.InvestmentPaymentSummaryByCompany(investmentPaymentService))
}
