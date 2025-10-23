package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/country"
)

func CountryRouter(app fiber.Router, svc country.Service) {
	app.Get("", handlers.GetAllCountries(svc))
	app.Get("/:id", handlers.GetCountryByID(svc))
}

func AdminCountryRouter(app fiber.Router, svc country.Service) {
	app.Get("", handlers.GetAllCountries(svc))
	app.Get("/:id", handlers.GetCountryByID(svc))
}

func CompanyCountryRouter(app fiber.Router, svc country.Service) {
	app.Get("", handlers.GetAllCountries(svc))
	app.Get("/:id", handlers.GetCountryByID(svc))
}
