package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// FindAllCompanies powers GET /api/companies. The frontend
// SelectCompany dropdown uses this to populate its options.
// Org-scope is applied by the service layer.
func FindAllCompanies(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CompanyFindAllRequest)
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
