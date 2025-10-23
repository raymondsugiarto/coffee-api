package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllCompanyCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		req.CompanyID = shared.GetUserCredential(c.Context()).CompanyID

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindCompanyCustomerByID(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByIDWithScope(c.Context(), id, []string{"complete"})
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
