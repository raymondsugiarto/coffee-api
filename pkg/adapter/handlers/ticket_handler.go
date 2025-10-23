package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/ticket"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateTicket(service ticket.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.TicketInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.UserID = shared.GetUserCredential(c.Context()).UserID

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func FindAllTicket(service ticket.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.FindAllRequest)
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

func FindTicketByID(service ticket.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateTicket(service ticket.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.TicketInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.ID = id
		dto.UserID = shared.GetUserCredential(c.Context()).UserID

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteTicket(service ticket.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
