package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	unitlink "github.com/raymondsugiarto/coffee-api/pkg/module/customer/unit_link"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
)

func GetCountCustomerByType(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.CountByType(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetCountNewCustomerByType(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.CountByTypeThisMonth(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UnitLinkSummaryByCompany(service unitlink.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.SummaryByCompany(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
func InvestmentDistributionSummaryByCompany(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.SummaryByCompany(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func InvestmentProductSummary(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.SummaryList(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UnitLinkSummaryPerType(service unitlink.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.SummaryPerType(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetCountCompanyByType(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.CountByType(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func InvestmentPaymentSummaryByCompany(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetPaymentSummary(c.Context())
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
