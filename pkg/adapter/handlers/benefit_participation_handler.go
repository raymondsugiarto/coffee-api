package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	bp "github.com/raymondsugiarto/coffee-api/pkg/module/benefit_participation"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateBenefitParticipation(service bp.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.BenefitParticipationDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		result, err := service.Create(c.Context(), request)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindBenefitParticipationByID(service bp.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
