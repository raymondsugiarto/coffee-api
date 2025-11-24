package middleware

import (
	"fmt"

	config "github.com/raymondsugiarto/coffee-api/config"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected protect routes
func Protected() fiber.Handler {
	cfg := config.GetConfig()
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(cfg.Server.Rest.SecretKey)},
		ErrorHandler:   jwtError,
		SuccessHandler: SuccessHandler,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func SuccessHandler(c *fiber.Ctx) error {
	token := c.Locals(entity.UserContextKey).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userCredentialsData := new(entity.UserCredentialData)
	userCredentialsData.ID = claims["id"].(string)
	userCredentialsData.UserID = claims["uid"].(string)

	if claims["aid"] != nil {
		userCredentialsData.AdminID = claims["aid"].(string)
	}
	if claims["coid"] != nil {
		userCredentialsData.CompanyID = claims["coid"].(string)
		c.Locals(entity.CompanyKey, userCredentialsData.CompanyID)
	}
	if claims["cid"] != nil {
		userCredentialsData.CustomerID = claims["cid"].(string)
	}
	fmt.Printf("%+v", userCredentialsData)
	c.Locals(entity.UserCredentialDataKey, userCredentialsData)
	return c.Next()
}
