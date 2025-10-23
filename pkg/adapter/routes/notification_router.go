package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/notification"
)

func CustomerNotificationRouter(app fiber.Router, svc notification.Service) {
	app.Get("", handlers.FindAllMyNotification(svc))
}
