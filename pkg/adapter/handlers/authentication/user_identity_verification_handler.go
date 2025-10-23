package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	useridentityverification "github.com/raymondsugiarto/coffee-api/pkg/module/user_identity_verification"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func ForgotPasswordCustomer(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.ForgotPasswordInputDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()

		response, err := service.ForgotPasswordCustomer(c.Context(), dto)
		if err != nil {
			log.Errorf("errorForgotPassword: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func ForgotPasswordCompany(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.ForgotPasswordInputDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()

		response, err := service.ForgotPasswordCompany(c.Context(), dto)
		if err != nil {
			log.Errorf("errorForgotPassword: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func VerifyUserIdentityVerificationForPassword(service useridentityverification.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.UserIdentityVerificationInputPasswordDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = c.Params("id")

		response, err := service.Verify(c.Context(), dto)
		if err != nil {
			log.Errorf("errorForgotPassword: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func VerifyUserIdentityVerificationForEmail(service useridentityverification.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.UserIdentityVerificationInputEmailDto)
		if err := c.BodyParser(&request); err != nil {
			fmt.Printf("error: %v\n", err)
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = c.Params("id")

		response, err := service.Verify(c.Context(), dto)
		if err != nil {
			log.Errorf("errVerifyEmail: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func ResendUserIdentityVerification(service useridentityverification.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		response, err := service.Resend(c.Context(), c.Params("id"))
		if err != nil {
			log.Errorf("errVerifyPhoneNumber: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}
