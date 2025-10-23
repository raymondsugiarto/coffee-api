package handlers

import (
	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	netassetvalue "github.com/raymondsugiarto/coffee-api/pkg/module/net_asset_value"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateNetAssetValue(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.NetAssetValueInputDto)

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

func CreateBatchNetAssetValue(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.NetAssetValueBatchInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		result, err := service.CreateBatch(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllNetAssetValue(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.NetAssetValueFindAllRequest)
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

func FindNetAssetValueByID(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateNetAssetValue(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.NetAssetValueInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id

		// TODO: remove previous file
		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteNetAssetValue(service netassetvalue.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
