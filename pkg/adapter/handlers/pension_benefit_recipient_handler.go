package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	pensionbenefitrecipient "github.com/raymondsugiarto/coffee-api/pkg/module/pension_benefit_recipient"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateBenefitRecipient(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.PensionBenefitRecipientInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		result, err := svc.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllBenefitRecipient(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.RecipientFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := svc.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindBenefitRecipientByID(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := svc.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateBenefitRecipient(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.PensionBenefitRecipientInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id

		result, err := svc.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteBenefitRecipient(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := svc.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func BacthDeleteBenefitRecipient(svc pensionbenefitrecipient.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customerID := c.Query("customer_id")
		request := new([]entity.PensionBenefitRecipientInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		recipientDtos := make([]*entity.PensionBenefitRecipientDto, len(*request))
		for i, recipient := range *request {
			recipientDtos[i] = recipient.ToDto()
		}

		err := svc.BatchDelete(c.Context(), customerID, recipientDtos)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
