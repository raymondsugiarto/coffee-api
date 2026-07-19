package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	salarycomponent "github.com/raymondsugiarto/coffee-api/pkg/module/salary_component"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllSalaryComponents(service salarycomponent.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.SalaryComponentFindAllRequest)
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

func FindOneSalaryComponent(service salarycomponent.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.Get(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func CreateSalaryComponent(service salarycomponent.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.SalaryComponentDto)
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

func UpdateSalaryComponent(service salarycomponent.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.SalaryComponentDto)
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

func DeleteSalaryComponent(service salarycomponent.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"deleted": id})
	}
}
