package routes

import (
	"github.com/gofiber/fiber/v2"
	h "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
)

func AdminRouter(app fiber.Router, svc admin.Service) {
	app.Get("", handlers.FindAllAdmin(svc))
}

func AdminMyRouter(app fiber.Router, svc admin.Service, permissionSvc permission.Service) {
	app.Get("", handlers.AdminGetMyProfile(svc))
	app.Get("/permissions", h.GetAdminMyPermission(permissionSvc))
	app.Post("/profile-picture", handlers.UpdateAdminProfileImage(svc))
	app.Put("/name", handlers.UpdateAdminName(svc))
}

func AdminCompanyRouter(app fiber.Router, svc admin.Service) {
	app.Post("", handlers.CreateAdminCompany(svc))
	app.Get("/:id", handlers.FindAllAdminByCompanyID(svc))
}
