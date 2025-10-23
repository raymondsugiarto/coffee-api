package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/admin"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
)

func AdminUnitLinkRouter(app fiber.Router, unitLinkService unitlink.Service) {
	app.Get("/unit-links", handlers.FindAllUnitLink(unitLinkService))
}
