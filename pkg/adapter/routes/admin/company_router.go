package routes

import (
	"github.com/gofiber/fiber/v2"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
)

func CompanyRouter(app fiber.Router, svc company.Service) {
	app.Post("", handlers.CreateCompany(svc))
	app.Get("", handlers.FindAllCompany(svc))
	app.Get("/:id", handlers.FindCompanyByID(svc))
	app.Delete("/:id", handlers.DeleteCompany(svc))
	app.Put("/:id", handlers.UpdateCompany(svc))
}
