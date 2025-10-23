package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/ticket"
)

func CustomerTicketRouter(app fiber.Router, svc ticket.Service) {
	app.Post("", handlers.CreateTicket(svc))
	app.Get("", handlers.FindAllTicket(svc))
	app.Get("/:id", handlers.FindTicketByID(svc))
	app.Put("/:id", handlers.UpdateTicket(svc))
	app.Delete("/:id", handlers.DeleteTicket(svc))
}

// func CustomerRewardRouter(app fiber.Router, svc reward.Service) {
// 	app.Get("", handlers.FindAllReward(svc))
// }
