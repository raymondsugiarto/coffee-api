package routes

import (
	"github.com/gofiber/fiber/v2"
	companyHandler "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/company"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
)

func CompanyCustomerRouter(app fiber.Router, svc customer.Service) {
	app.Get("", companyHandler.FindAllCompanyCustomer(svc))
	app.Get("/:id", companyHandler.FindCompanyCustomerByID(svc))
	app.Post("", handlers.CreateCustomer(svc))
	app.Delete("/:id", handlers.DeleteCustomer(svc))
	app.Put("/:id", handlers.UpdateCustomer(svc))
	app.Post("/:id/suspend", handlers.SuspendStatusCustomer(svc))
}
