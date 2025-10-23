package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
	rolepermission "github.com/raymondsugiarto/coffee-api/pkg/module/permission/role_permission"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
)

func GetAdminMyPermission(service permission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID

		fmt.Printf("ud: %v", shared.GetUserCredential(c.Context()))
		result, err := service.FindAll(c.Context(), &entity.PermissionFindAllRequest{
			UserID: id,
		})
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func AddPermissionToRole(service rolepermission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID := c.Params("id")
		permissionID := c.Params("permissionId")

		if err := service.AddPermissionToRole(c.Context(), roleID, permissionID); err != nil {
			return err
		}

		return c.JSON("success")
	}
}

func RemovePermissionFromRole(service rolepermission.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID := c.Params("id")
		permissionID := c.Params("permissionId")

		if err := service.RemovePermissionFromRole(c.Context(), roleID, permissionID); err != nil {
			return err
		}

		return c.JSON("success")
	}
}
