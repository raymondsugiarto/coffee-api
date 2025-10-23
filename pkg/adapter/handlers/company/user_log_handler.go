package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	userlog "github.com/raymondsugiarto/coffee-api/pkg/module/user_log"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllCompanyUserLog(service userlog.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.UserLogFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		req.CompanyID = shared.GetCompanyID(c.Context())
		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
