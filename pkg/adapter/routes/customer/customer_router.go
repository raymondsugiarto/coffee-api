package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
)

func AdminCustomerRouter(app fiber.Router, svc customer.Service) {
	app.Post("", handlers.CreateCustomer(svc))
	app.Get("", handlers.FindAllCustomer(svc))
	app.Get("/:id", handlers.FindCustomerByID(svc))
	app.Delete("/:id", handlers.DeleteCustomer(svc))
	app.Put("/:id", handlers.UpdateCustomer(svc))
	app.Post("/:id/suspend", handlers.SuspendStatusCustomer(svc))
}

func CustomerRouter(app fiber.Router, svc customer.Service) {
	app.Get("", handlers.FindAllMyReferral(svc))
}
