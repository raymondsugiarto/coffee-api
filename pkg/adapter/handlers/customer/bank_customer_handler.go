package customer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	bankcustomer "github.com/raymondsugiarto/coffee-api/pkg/module/customer/bank_customer"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateBankCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.BankCustomerInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindBankCustomerByID(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateBankCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.BankCustomerInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.ID = c.Params("id")

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteBankCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func FindAllBankCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.BankCustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindByCustomer(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindBankCustomerByCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.BankCustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindByCustomer(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func SetDefaultBankCustomer(service bankcustomer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.SetDefaultBankCustomer(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
