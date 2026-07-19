package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	cashdebt "github.com/raymondsugiarto/coffee-api/pkg/module/cash_debt"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// FindAllCashDebts powers GET /api/cash-debts.
func FindAllCashDebts(service cashdebt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CashDebtFindAllRequest)
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

// FindOneCashDebt powers GET /api/cash-debts/:id.
func FindOneCashDebt(service cashdebt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.Get(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// CreateCashDebt powers POST /api/cash-debts.
func CreateCashDebt(service cashdebt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CashDebtDto)
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

// UpdateCashDebt powers PUT /api/cash-debts/:id.
func UpdateCashDebt(service cashdebt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.CashDebtDto)
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

// DeleteCashDebt powers DELETE /api/cash-debts/:id.
func DeleteCashDebt(service cashdebt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"deleted": id})
	}
}
