package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/item"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllItems(service item.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		itemReq := new(entity.ItemFindAllRequest)
		if err := c.QueryParser(itemReq); err != nil {
			return status.New(status.BadRequest, err)
		}

		userCred := shared.GetUserCredential(c.Context())
		itemReq.UserID = userCred.UserID

		result, err := service.FindAll(c.Context(), itemReq)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
