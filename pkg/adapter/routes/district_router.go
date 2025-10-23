package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/district"
)

func DistrictRouter(app fiber.Router, svc district.Service) {
	app.Get("", handlers.GetAllDistricts(svc))
	app.Get("/:id", handlers.GetDistrictByID(svc))
}

func AdminDistrictRouter(app fiber.Router, svc district.Service) {
	app.Get("", handlers.GetAllDistricts(svc))
	app.Get("/:id", handlers.GetDistrictByID(svc))
}

func CompanyDistrictRouter(app fiber.Router, svc district.Service) {
	app.Get("", handlers.GetAllDistricts(svc))
	app.Get("/:id", handlers.GetDistrictByID(svc))
}
