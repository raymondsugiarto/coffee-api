package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/claim"
)

func CustomerClaimRouter(app fiber.Router, svc claim.Service) {
	app.Post("", handlers.CreateClaim(svc))
	app.Get("", handlers.FindAllClaimByCustomer(svc))
}

func CompanyClaimRouter(app fiber.Router, svc claim.Service) {
	app.Post("", handlers.CreateClaim(svc))
	app.Get("", handlers.FindClaimByCompany(svc))
	app.Get("/:id", handlers.FindClaimByID(svc))
}
