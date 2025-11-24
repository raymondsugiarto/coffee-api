package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func SignIn(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.LoginRequestDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		response, err := service.SignIn(c.Context(), request)
		if err != nil {
			fmt.Printf("errorSignIn: %v\n", err)
			return err
		}

		return c.JSON(response)
	}
}
