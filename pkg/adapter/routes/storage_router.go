package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
)

func StorageFileRouter(app fiber.Router) {
	app.Get("", handlers.GetStorageFile())
}
