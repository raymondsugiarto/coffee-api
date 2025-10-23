package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/estatement"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func GenerateEstatement(service estatement.EStatementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.EstatementRequestDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		customerId := shared.GetUserCredential(c.Context()).CustomerID
		req.CustomerID = customerId

		result, err := service.GenerateEstatement(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GenerateAndSendEstatementEmail(service estatement.EStatementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.EstatementEmailRequestDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		customerId := shared.GetUserCredential(c.Context()).CustomerID
		req.CustomerID = customerId

		go service.GenerateAndSendEmail(c.Context(), req)

		return c.JSON(fiber.Map{
			"message": "E-statement email will be sent shortly",
			"status":  "OK",
		})
	}
}
