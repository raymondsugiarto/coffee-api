package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	pensionbenefitrecipient "github.com/raymondsugiarto/coffee-api/pkg/module/pension_benefit_recipient"
)

func AdminBenefitRecipientRouter(app fiber.Router, svc pensionbenefitrecipient.Service) {
	app.Post("", handlers.CreateBenefitRecipient(svc))
	app.Get("", handlers.FindAllBenefitRecipient(svc))
	app.Get("/:id", handlers.FindBenefitRecipientByID(svc))
	app.Put("/:id", handlers.UpdateBenefitRecipient(svc))
	app.Delete("/:id", handlers.DeleteBenefitRecipient(svc))
	app.Delete("*", handlers.BacthDeleteBenefitRecipient(svc))
}

func CompanyBenefitRecipientRouter(app fiber.Router, svc pensionbenefitrecipient.Service) {
	app.Post("", handlers.CreateBenefitRecipient(svc))
	app.Get("", handlers.FindAllBenefitRecipient(svc))
	app.Get("/:id", handlers.FindBenefitRecipientByID(svc))
	app.Put("/:id", handlers.UpdateBenefitRecipient(svc))
	app.Delete("/:id", handlers.DeleteBenefitRecipient(svc))
	app.Delete("*", handlers.BacthDeleteBenefitRecipient(svc))
}

func CustomerBenefitRecipientRouter(app fiber.Router, svc pensionbenefitrecipient.Service) {
	app.Post("", handlers.CreateBenefitRecipient(svc))
	app.Get("", handlers.FindAllBenefitRecipient(svc))
	app.Get("/:id", handlers.FindBenefitRecipientByID(svc))
	app.Put("/:id", handlers.UpdateBenefitRecipient(svc))
	app.Delete("/:id", handlers.DeleteBenefitRecipient(svc))
	app.Delete("*", handlers.BacthDeleteBenefitRecipient(svc))
}
