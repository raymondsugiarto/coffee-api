package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/bank"
)

func AdminBankRouter(app fiber.Router, svc bank.Service) {
	app.Post("", handlers.CreateBank(svc))
	app.Get("", handlers.FindAllBank(svc))
	app.Get("/:id", handlers.FindBankByID(svc))
	app.Delete("/:id", handlers.DeleteBank(svc))
	app.Put("/:id", handlers.UpdateBank(svc))
}

func CustomerBankRouter(app fiber.Router, svc bank.Service) {
	app.Get("", handlers.FindAllBank(svc))
	app.Get("/:id", handlers.FindBankByID(svc))
}
