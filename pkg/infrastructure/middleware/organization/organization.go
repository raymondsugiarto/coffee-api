package organization

import (
	"errors"

	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/database"
	"github.com/raymondsugiarto/coffee-api/pkg/model"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	return func(c *fiber.Ctx) error {

		// Get id from request, else we generate one
		origin := c.Get(cfg.HeaderOriginKey)
		if origin == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed origin", "data": nil})
		}

		originType := c.Get(cfg.HeaderOriginTypeKey)
		if originType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed origin", "data": nil})
		}

		log.Infof("Origin: %v", origin)
		org, err := getOrganizationByOrigin(origin)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Organization not found", "data": nil})
		}
		c.Locals(entity.OriginKey, origin)
		c.Locals(entity.OriginTypeKey, originType)

		c.Locals(entity.OrganizationKey, &entity.OrganizationData{
			ID: org.ID,
		})

		return c.Next()
	}
}

func getOrganizationByOrigin(origin string) (*model.Organization, error) {
	db := database.DBConn

	var organization model.Organization
	err := db.Where("origin = ?", origin).Find(&organization).Error
	if err != nil {
		log.Errorf("Error: %v", err)
		return nil, errors.New("organization not found")
	}
	if organization.ID == "" {
		return nil, errors.New("organization not found")
	}

	return &organization, nil
}
