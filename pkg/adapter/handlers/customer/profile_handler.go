package customer

import (
	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/profile"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CustomerGetMyProfile(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID

		result, err := service.FindByUserID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateMyProfile(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).CustomerID
		request := new(entity.UpdateCustomerKYCDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id
		dto.PensionBenefitRecipients[0].CustomerID = id
		customerAttachment(c, dto)

		result, err := service.UpdateMyProfile(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
