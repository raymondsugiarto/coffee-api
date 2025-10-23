package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
)

func InvestmentProductRouter(app fiber.Router, svc investmentproduct.Service) {
	app.Post("", handlers.CreateInvestmentProduct(svc))
	app.Get("", handlers.FindAllInvestmentProduct(svc))
	app.Get("/:id", handlers.FindInvestmentProductByID(svc))
	app.Delete("/:id", handlers.DeleteInvestmentProduct(svc))
	app.Put("/:id", handlers.UpdateInvestmentProduct(svc))
}

func CustomerInvestmentProductRouter(app fiber.Router, svc investmentproduct.Service) {
	app.Get("", handlers.FindAllInvestmentProduct(svc))
	app.Get("/:id", handlers.FindInvestmentProductByID(svc))
}
