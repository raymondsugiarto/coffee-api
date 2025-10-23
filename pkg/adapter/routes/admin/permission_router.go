package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
)

func PermissionRouter(app fiber.Router, permissionSvc permission.Service) {
	app.Get("", handlers.FindAllPermission(permissionSvc))
}
