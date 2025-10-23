package sharedorganization

import (
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/pagination"

	"github.com/gofiber/fiber/v2"
)

func AddFilterOrganizationData(c *fiber.Ctx, req pagination.PaginationRequestDto) {
	req.AddFilter(pagination.FilterItem{
		Field: "organization_id",
		Val:   c.Locals(entity.OrganizationKey).(*entity.OrganizationData).ID,
		Op:    "eq",
	})
}
