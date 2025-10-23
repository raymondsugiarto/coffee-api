package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/reward"
)

func AdminRewardRouter(app fiber.Router, svc reward.Service) {
	app.Post("", handlers.CreateReward(svc))
	app.Get("", handlers.FindAllReward(svc))
	app.Get("/:id", handlers.FindRewardByID(svc))
	app.Put("/:id", handlers.UpdateReward(svc))
	app.Delete("/:id", handlers.DeleteReward(svc))
}

func CustomerRewardRouter(app fiber.Router, svc reward.Service) {
	app.Get("", handlers.FindAllReward(svc))
}
