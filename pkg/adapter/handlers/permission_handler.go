package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllPermission(service permission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.PermissionFindAllRequest)
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
