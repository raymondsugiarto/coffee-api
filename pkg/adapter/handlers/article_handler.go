package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/article"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.ArticleInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		if err := articleAttachment(c, dto); err != nil {
			return err
		}

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func FindAllArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ArticleFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindArticleByID(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindArticleBySlug(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")

		result, err := service.FindBySlug(c.Context(), slug)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.ArticleInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.ID = id
		if err := articleAttachment(c, dto); err != nil {
			return err
		}

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateArticleStatus(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.ArticleUpdateStatusDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		id := c.Params("id")
		request.ID = id

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func DeleteArticle(service article.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func articleAttachment(c *fiber.Ctx, dto *entity.ArticleDto) error {
	articleImageFile, errArticleImageFile := c.FormFile("imageFile")

	if articleImageFile != nil {
		articleImageFilePath := fmt.Sprintf("./storage/article/%s", articleImageFile.Filename)
		if err := c.SaveFile(articleImageFile, articleImageFilePath); err != nil {
			return status.New(status.BadRequest, errArticleImageFile)
		}
		dto.ImageUrl = articleImageFilePath
	}

	if dto.ImageUrl == "" {
		return status.New(status.BadRequest, fmt.Errorf("image is required"))
	}
	return nil
}
