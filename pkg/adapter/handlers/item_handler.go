package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
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

func FindOneItem(service item.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func CreateItem(service item.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ItemDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		result, err := service.Create(c.Context(), req)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func UpdateItem(service item.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.ItemDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		req.ID = id
		result, err := service.Update(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func DeleteItem(service item.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"deleted": id})
	}
}
