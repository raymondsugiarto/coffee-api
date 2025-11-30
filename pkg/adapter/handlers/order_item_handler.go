package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	orderitem "github.com/raymondsugiarto/coffee-api/pkg/module/order/order_item"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CountMyOrderItems(service orderitem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemReq := new(entity.OrderFindAllRequest)
		if err := c.QueryParser(itemReq); err != nil {
			return status.New(status.BadRequest, err)
		}

		userCred := shared.GetUserCredential(c.Context())
		itemReq.AdminID = userCred.AdminID

		result, err := service.Count(c.Context(), itemReq)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
