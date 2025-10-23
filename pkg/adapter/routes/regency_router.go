package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/regency"
)

func RegencyRouter(app fiber.Router, svc regency.Service) {
	app.Get("", handlers.GetAllRegencies(svc))
	app.Get("/:id", handlers.GetRegencyByID(svc))
}

func AdminRegencyRouter(app fiber.Router, svc regency.Service) {
	app.Get("", handlers.GetAllRegencies(svc))
	app.Get("/:id", handlers.GetRegencyByID(svc))
}

func CompanyRegencyRouter(app fiber.Router, svc regency.Service) {
	app.Get("", handlers.GetAllRegencies(svc))
	app.Get("/:id", handlers.GetRegencyByID(svc))
}
