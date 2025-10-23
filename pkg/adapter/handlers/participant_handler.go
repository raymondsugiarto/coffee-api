package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer/participant"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllMyParticipant(service participant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ParticipantFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		req.CustomerID = shared.GetUserCredential(c.Context()).CustomerID
		req.Status = model.ParticipantStatusActive

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindParticipantByID(service participant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllParticipantByCompany(service participant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ParticipantFindAllRequest)
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
