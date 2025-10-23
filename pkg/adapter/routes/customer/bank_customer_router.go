package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/customer"
	bankcustomer "github.com/raymondsugiarto/coffee-api/pkg/module/customer/bank_customer"
)

func CustomerBankCustomerRouter(app fiber.Router, svc bankcustomer.Service) {
	app.Post("", handlers.CreateBankCustomer(svc))
	app.Get("/:id", handlers.FindBankCustomerByID(svc))
	app.Put("/:id", handlers.UpdateBankCustomer(svc))
	app.Patch("/:id/default", handlers.SetDefaultBankCustomer(svc))
	app.Delete("/:id", handlers.DeleteBankCustomer(svc))
	app.Get("", handlers.FindAllBankCustomer(svc))
}
