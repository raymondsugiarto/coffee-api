package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/driver"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllDrivers(service driver.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.DriverFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		result, err := service.FindAllDrivers(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
