package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
)

func CustomerParticipantRouter(app fiber.Router, svc participant.Service) {
	app.Get("", handlers.FindAllMyParticipant(svc))
	app.Get("/:id", handlers.FindParticipantByID(svc))
}

func CompanyParticipantRouter(app fiber.Router, svc participant.Service) {
	app.Get("", handlers.FindAllParticipantByCompany(svc))
	app.Get("/:id", handlers.FindParticipantByID(svc))
}
