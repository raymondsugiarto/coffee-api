package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	itemcategory "github.com/raymondsugiarto/coffee-api/pkg/module/item_category"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllItemCategories(service itemcategory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ItemCategoryFindAllRequest)
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

func FindOneItemCategory(service itemcategory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.Get(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func CreateItemCategory(service itemcategory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ItemCategoryDto)
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

func UpdateItemCategory(service itemcategory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.ItemCategoryDto)
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

func DeleteItemCategory(service itemcategory.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"deleted": id})
	}
}
