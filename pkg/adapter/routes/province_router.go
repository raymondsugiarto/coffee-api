package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/province"
)

func ProvinceRouter(app fiber.Router, svc province.Service) {
	app.Get("", handlers.GetAllProvinces(svc))
	app.Get("/:id", handlers.GetProvinceByID(svc))
}

func AdminProvinceRouter(app fiber.Router, svc province.Service) {
	app.Get("", handlers.GetAllProvinces(svc))
	app.Get("/:id", handlers.GetProvinceByID(svc))
}

func CompanyProvinceRouter(app fiber.Router, svc province.Service) {
	app.Get("", handlers.GetAllProvinces(svc))
	app.Get("/:id", handlers.GetProvinceByID(svc))
}
