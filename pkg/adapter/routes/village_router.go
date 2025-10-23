package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/village"
)

func VillageRouter(app fiber.Router, svc village.Service) {
	app.Get("", handlers.GetAllVillages(svc))
	app.Get("/:id", handlers.GetVillageByID(svc))
}

func AdminVillageRouter(app fiber.Router, svc village.Service) {
	app.Get("", handlers.GetAllVillages(svc))
	app.Get("/:id", handlers.GetVillageByID(svc))
}

func CompanyVillageRouter(app fiber.Router, svc village.Service) {
	app.Get("", handlers.GetAllVillages(svc))
	app.Get("/:id", handlers.GetVillageByID(svc))
}
