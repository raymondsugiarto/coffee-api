package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	netassetvalue "github.com/raymondsugiarto/coffee-api/pkg/module/net_asset_value"
)

func NetAssetValueRouter(app fiber.Router, svc netassetvalue.Service) {
	app.Post("", handlers.CreateNetAssetValue(svc))
	app.Post("/batches", handlers.CreateBatchNetAssetValue(svc))
	app.Get("", handlers.FindAllNetAssetValue(svc))
	app.Get("/:id", handlers.FindNetAssetValueByID(svc))
	app.Delete("/:id", handlers.DeleteNetAssetValue(svc))
	app.Put("/:id", handlers.UpdateNetAssetValue(svc))
}
