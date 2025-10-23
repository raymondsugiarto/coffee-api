package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
)

func AdminUploadRouter(app fiber.Router, svc customer.Service) {
	app.Post("/customers", handlers.UploadCustomer(svc))
}

func CompanyUploadRouter(app fiber.Router, svc customer.Service) {
	app.Post("/customers", handlers.UploadCustomer(svc))
}
