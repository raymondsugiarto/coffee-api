package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/redeem"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateRedeem(service redeem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.RedeemInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		customerID := shared.GetUserCredential(c.Context()).CustomerID
		OrganizationID := shared.GetOrganization(c.Context()).ID
		dto.OrganizationID = OrganizationID
		dto.CustomerID = customerID

		result, err := service.RedeemReward(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func GetCustomerRedeems(service redeem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.RedeemFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		// Set customer ID from authenticated user context
		customerID := shared.GetUserCredential(c.Context()).CustomerID
		req.CustomerID = customerID

		result, err := service.GetAllRedemptions(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetAllRedeems(service redeem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.RedeemFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.GetAllRedemptions(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetRedeemByID(service redeem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.GetRedemptionByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateRedeemStatus(service redeem.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.UpdateRedeemStatusDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		result, err := service.UpdateRedemptionStatus(c.Context(), id, request)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
