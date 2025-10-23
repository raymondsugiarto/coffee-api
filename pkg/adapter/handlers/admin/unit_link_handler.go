package handlers

import (
	"github.com/gofiber/fiber/v2"
	ec "github.com/raymondsugiarto/coffee-api/pkg/entity/customer"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllUnitLink(service unitlink.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(ec.UnitLinkFindAllRequest)
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

func SumInvestmentProductByCustomer(service unitlink.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.SumInvestmentProductByCustomer(
			c.Context(),
			shared.GetUserCredential(c.Context()).CustomerID,
		)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllPortfolioWithNav(service unitlink.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(ec.PortfolioFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		result, err := service.FindAllPortfolioWithNav(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
