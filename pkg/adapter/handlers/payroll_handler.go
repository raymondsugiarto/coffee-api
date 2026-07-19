package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/payroll"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// SimulatePayroll powers POST /api/payroll/simulate. Returns the
// per-session breakdown + rolled totals without persisting
// anything. Read-only — safe to call repeatedly.
func SimulatePayroll(service payroll.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.SimulatePayrollRequest)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		result, err := service.Simulate(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// SavePayroll powers POST /api/payroll. Persists the header +
// components the operator just approved in the simulation step.
func SavePayroll(service payroll.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.SavePayrollRequest)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		result, err := service.Save(c.Context(), req)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

// FindAllPayrolls powers GET /api/payroll.
func FindAllPayrolls(service payroll.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.EmployeeSalaryFindAllRequest)
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

// FindOnePayroll powers GET /api/payroll/:id.
func FindOnePayroll(service payroll.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.FindOne(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
