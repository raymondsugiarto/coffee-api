package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	investmentproduct "github.com/raymondsugiarto/coffee-api/pkg/module/investment_product"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateInvestmentProduct(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentProductInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		investmentProductAttachment(c, dto)

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllInvestmentProduct(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentProductFilter)
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

func FindInvestmentProductByID(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.InvestmentProductFilter)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindByID(c.Context(), id, req.IncludeAum)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateInvestmentProduct(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.InvestmentProductInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id
		investmentProductAttachment(c, dto)

		// TODO: remove previous file
		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteInvestmentProduct(service investmentproduct.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func investmentProductAttachment(c *fiber.Ctx, dto *entity.InvestmentProductDto) error {
	fundFactSheet, errFundFactSheet := c.FormFile("fundFactSheet")
	riplay, errRiplay := c.FormFile("riplay")

	// upload file if exist
	if fundFactSheet != nil {
		fundFactSheetPath := fmt.Sprintf("./storage/%s", fundFactSheet.Filename)
		// Save the files to disk:
		if err := c.SaveFile(fundFactSheet, fundFactSheetPath); err != nil {
			return status.New(status.BadRequest, errFundFactSheet)
		}
		dto.FundFactSheet = fundFactSheetPath
	}

	// upload file if exist
	if riplay != nil {
		riplayPath := fmt.Sprintf("./storage/%s", riplay.Filename)
		// Save the files to disk:
		if err := c.SaveFile(riplay, riplayPath); err != nil {
			return status.New(status.BadRequest, errRiplay)
		}
		dto.Riplay = riplayPath
	}
	return nil
}
