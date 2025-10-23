package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	bp "github.com/raymondsugiarto/coffee-api/pkg/module/benefit_participation"
)

func BenefitParticipationRouter(app fiber.Router, svc bp.Service) {
	app.Post("", handlers.CreateBenefitParticipation(svc))
	app.Get("/:id", handlers.FindBenefitParticipationByID(svc))
}
