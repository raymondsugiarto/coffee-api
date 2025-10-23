package routes

import (
	"github.com/gofiber/fiber/v2"
	h "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/permission/permission"
	userlog "github.com/raymondsugiarto/coffee-api/pkg/module/user_log"
)

func CompanyMyRouter(app fiber.Router, svc company.Service, permissionSvc permission.Service) {
	app.Get("", handlers.CompanyGetMyProfile(svc))
	app.Get("/permissions", h.GetAdminMyPermission(permissionSvc))
	app.Get("/profile/:id", handlers.FindCompanyByID(svc))
}

func FindAllCompanyUserLog(app fiber.Router, svc userlog.Service) {
	app.Get("", handlers.FindAllCompanyUserLog(svc))
}
