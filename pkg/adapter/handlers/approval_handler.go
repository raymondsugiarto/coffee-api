package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindApprovalByID(svc approval.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := svc.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateApproval(svc approval.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.ApprovalDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request
		dto.ID = id

		result, err := svc.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteApproval(svc approval.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := svc.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func FindAllApproval(svc approval.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ApprovalFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := svc.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func ConfirmationApproval(svc approval.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.ApprovalDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request
		dto.ID = id

		result, err := svc.Confirmation(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
