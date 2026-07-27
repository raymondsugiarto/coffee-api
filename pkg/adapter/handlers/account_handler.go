package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	accounting "github.com/raymondsugiarto/coffee-api/pkg/module/accounting"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// FindAllAccounts powers GET /api/accounts.
func FindAllAccounts(service accounting.AccountService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.AccountFindAllRequest)
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

// FindOneAccount powers GET /api/accounts/:id.
func FindOneAccount(service accounting.AccountService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.Get(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// CreateAccount powers POST /api/accounts.
func CreateAccount(service accounting.AccountService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.AccountDto)
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

// UpdateAccount powers PUT /api/accounts/:id.
func UpdateAccount(service accounting.AccountService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.AccountDto)
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

// DeleteAccount powers DELETE /api/accounts/:id.
func DeleteAccount(service accounting.AccountService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"deleted": id})
	}
}
