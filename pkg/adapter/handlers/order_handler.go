package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/order"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateOrder(service order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.OrderInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		userCred := shared.GetUserCredential(c.Context())
		dto := request.ToDto()
		dto.AdminID = userCred.AdminID

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func FindAllMyOrders(service order.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemReq := new(entity.OrderFindAllRequest)
		if err := c.QueryParser(itemReq); err != nil {
			return status.New(status.BadRequest, err)
		}

		userCred := shared.GetUserCredential(c.Context())
		itemReq.AdminID = userCred.AdminID

		result, err := service.FindAll(c.Context(), itemReq)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
