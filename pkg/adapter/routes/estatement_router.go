package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/estatement"
)

func EstatementRouter(app fiber.Router, service estatement.EStatementService) {
	app.Post("", handlers.GenerateEstatement(service))
	app.Post("/send-email", handlers.GenerateAndSendEstatementEmail(service))
}
