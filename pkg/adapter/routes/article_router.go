package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/article"
)

func AdminArticleRouter(app fiber.Router, svc article.Service) {
	app.Post("", handlers.CreateArticle(svc))
	app.Get("", handlers.FindAllArticle(svc))
	app.Get("/:id", handlers.FindArticleByID(svc))
	app.Put("/:id", handlers.UpdateArticle(svc))
	app.Patch("/:id/status", handlers.UpdateArticleStatus(svc))
	app.Delete("/:id", handlers.DeleteArticle(svc))
}

func CustomerArticleRouter(app fiber.Router, svc article.Service) {
	app.Get("", handlers.FindAllArticle(svc))
	app.Get("/:slug", handlers.FindArticleBySlug(svc))
}
