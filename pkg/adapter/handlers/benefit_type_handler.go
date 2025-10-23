package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/benefit_type"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateBenefitType(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.CreateBenefitTypeRequest)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		result, err := service.Create(c.Context(), request.ToDto())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllBenefitType(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.BenefitTypeFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindBenefitTypeByID(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateBenefitType(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.UpdateBenefitTypeRequest)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		result, err := service.Update(c.Context(), id, request.ToDto())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteBenefitType(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func FindAllActiveBenefitType(service benefit_type.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.BenefitTypeFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		// Force status to ACTIVE for customer access
		req.Status = "ACTIVE"

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
