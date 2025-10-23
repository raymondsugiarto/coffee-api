package middleware

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response"
)

func DefaultErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		log.WithContext(c.Context()).Errorf("Error: %+v", err)
		resp := response.NewError(err)
		log.WithContext(c.Context()).Errorf("HTTP Code: %+v", resp.HTTPCode)
		return c.Status(resp.HTTPCode).JSON(resp)
	}

}

func DefaultResponseHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			return err
		}

		contentType := c.GetRespHeader("Content-Type")
		if contentType != "application/json" {
			return nil
		}
		body := c.Response().Body()

		output := response.NewSuccess(c.Response().StatusCode(), body)

		if len(body) > 0 {
			var data any
			if err := json.Unmarshal(body, &data); err == nil {
				output = response.NewSuccess(c.Response().StatusCode(), data)
			}
		}

		return c.JSON(output)
	}
}
