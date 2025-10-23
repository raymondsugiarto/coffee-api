package handlers

import (
	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	investmentitem "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_item"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllCompanyInvestmentItem(service investmentitem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentItemFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		id := c.Params("investmentID")
		if id != "" {
			req.InvestmentID = id
		}

		req.ShowAll = true
		req.CompanyID = shared.GetUserCredential(c.Context()).CompanyID

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllMyInvestmentItem(service investmentitem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentItemFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		req.ShowAll = true
		req.CustomerID = shared.GetUserCredential(c.Context()).CustomerID

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindInvestmentItemByID(service investmentitem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
