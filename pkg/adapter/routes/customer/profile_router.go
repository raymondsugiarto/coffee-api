package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
)

func ProfileRouter(app fiber.Router, customerService customer.Service) {
	// TODO: profile employee
	app.Get("/me", handlers.CustomerGetMyProfile(customerService))
	app.Put("/kyc", handlers.UpdateMyProfile(customerService))
	app.Post("/change-password", handlers.ChangePassword(customerService))
}

func UnitLinkRouter(app fiber.Router, unitLinkService unitlink.Service) {
	app.Get("/unit-links", handlers.FindAllPortolio(unitLinkService))

	// TODO: need to be refactor
	app.Get("/portfolio", handlers.FindAllPortolio(unitLinkService))
	app.Get("/portfolio/summary", handlers.SumInvestmentProductByCustomer(unitLinkService))
}
