package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		imagePath, err := handleFileUpload(c, "confirmationImage")
		if err != nil {
			return status.New(status.BadRequest, err)
		}
		if imagePath != "" {
			c.Locals("confirmationImage", imagePath)
		}

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func CreateInvestmentFromPaymentFee(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		imagePath, err := handleFileUpload(c, "confirmationImage")
		if err != nil {
			return status.New(status.BadRequest, err)
		}
		if imagePath != "" {
			c.Locals("confirmationImage", imagePath)
		}

		result, err := service.CreateFromPaymentFee(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentFindAllRequest)
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

func FindAllMyInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		req.CustomerID = shared.GetUserCredential(c.Context()).CustomerID
		req.IncludePayments = true
		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindInvestmentByID(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.InvestmentInputDto)
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

func DeleteInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func FindAllCompanyInvestment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		req.ShowAll = true
		req.CompanyID = &shared.GetUserCredential(c.Context()).CompanyID

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func handleFileUpload(c *fiber.Ctx, formField string) (string, error) {
	file, err := c.FormFile(formField)
	if file != nil && err == nil {
		filePath := fmt.Sprintf("./storage/%s", file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			return "", err
		}
		return filePath, nil
	}
	return "", nil
}
