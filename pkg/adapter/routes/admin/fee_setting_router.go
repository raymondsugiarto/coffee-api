package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/fee_setting"
)

func FeeSettingRouter(app fiber.Router, svc fee_setting.Service) {
	app.Post("", handlers.UpsertFeeSetting(svc))
	app.Get("", handlers.GetFeeSetting(svc))
}
