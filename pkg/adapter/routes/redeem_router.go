package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/redeem"
)

func CustomerRedeemRouter(app fiber.Router, svc redeem.Service) {
	app.Post("", handlers.CreateRedeem(svc))
	app.Get("", handlers.GetCustomerRedeems(svc))
	app.Get("/:id", handlers.GetRedeemByID(svc))
}

func AdminRedeemRouter(app fiber.Router, svc redeem.Service) {
	app.Get("", handlers.GetAllRedeems(svc))
	app.Get("/:id", handlers.GetRedeemByID(svc))
	app.Put("/:id/status", handlers.UpdateRedeemStatus(svc))
}
