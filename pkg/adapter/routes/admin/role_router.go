package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	rolepermission "github.com/raymondsugiarto/coffee-api/pkg/module/permission/role_permission"
	"github.com/raymondsugiarto/coffee-api/pkg/module/role"
)

func RoleRouter(app fiber.Router, svc role.Service, rolePermissionSvc rolepermission.Service) {
	app.Post("", handlers.CreateRole(svc))
	app.Get("", handlers.FindAllRole(svc))
	app.Get("/:id", handlers.FindRoleByID(svc))
	app.Delete("/:id", handlers.DeleteRole(svc))
	app.Put("/:id", handlers.UpdateRole(svc))

	app.Post("/:id/permissions/:permissionId", handlers.AddPermissionToRole(rolePermissionSvc))
	app.Delete("/:id/permissions/:permissionId", handlers.RemovePermissionFromRole(rolePermissionSvc))
}
