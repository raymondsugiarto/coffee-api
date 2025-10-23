package handlers

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/investment"
	investmentpayment "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_payment"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateInvestmentPayment(service investment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentPaymentInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.InvestmentID = c.Params("id")

		result, err := service.UploadPayment(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllInvestmentPayment(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentPaymentFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
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

func FindAllMyInvestmentPayment(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentPaymentFindAllRequest)
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

func FindInvestmentPaymentByID(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateInvestmentPaymentConfirmation(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		paymentID := c.Params("paymentId")
		request := new(entity.InvestmentPaymentConfirmationInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = paymentID

		cooperationAgreement, errCooperationAgreement := c.FormFile("confirmationImage")
		if cooperationAgreement != nil {
			cooperationAgreementPath := fmt.Sprintf("./storage/%s", cooperationAgreement.Filename)
			// Save the files to disk:
			if err := c.SaveFile(cooperationAgreement, cooperationAgreementPath); err != nil {
				return status.New(status.BadRequest, errCooperationAgreement)
			}
			dto.ConfirmationImageUrl = cooperationAgreementPath
		}

		result, err := service.UpdatePaymentConfirmation(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteInvestmentPayment(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func GetTotalMonthlyCompanyContribution(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.GetTotalMonthlyCompanyContributionRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.GetTotalMonthlyCompanyContribution(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllAdminInvestmentPayment(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentPaymentFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateExcel(c.Context(), req, result)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", "attachment; filename=investment_payments.xlsx")
			return c.SendStream(bytes.NewReader(excelBytes))
		}

		return c.JSON(result)
	}
}

func FindAdminInvestmentPaymentByID(service investmentpayment.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
