package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/province"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func GetAllProvinces(service province.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ProvinceFindAllRequest)
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

func GetProvinceByID(service province.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
