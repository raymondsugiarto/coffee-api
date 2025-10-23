package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func UpsertFeeSetting(service fee_setting.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.FeeSettingDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		err := service.UpsertConfig(c.Context(), request)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetFeeSetting(service fee_setting.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetConfig(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
