package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
)

func AdminInvestmentPaymentRouter(app fiber.Router, paymentService investmentpayment.Service) {
	app.Get("/investments/payments", handlers.FindAllAdminInvestmentPayment(paymentService))
	app.Get("/investments/payments/:id", handlers.FindAdminInvestmentPaymentByID(paymentService))
}
