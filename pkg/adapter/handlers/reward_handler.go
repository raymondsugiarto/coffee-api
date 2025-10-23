package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/reward"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateReward(service reward.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.RewardInputDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		if err := rewardAttachment(c, dto); err != nil {
			return err
		}

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

func FindAllReward(service reward.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.FindAllRequest)
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

func FindRewardByID(service reward.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateReward(service reward.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.RewardInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		dto.ID = id

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteReward(service reward.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func rewardAttachment(c *fiber.Ctx, dto *entity.RewardDto) error {
	rewardImage, errRewardImage := c.FormFile("imageFile")

	if rewardImage != nil {
		rewardImagePath := "./storage/" + rewardImage.Filename

		if err := c.SaveFile(rewardImage, rewardImagePath); err != nil {
			return status.New(status.BadRequest, errRewardImage)
		}
		dto.ImageUrl = rewardImagePath
	}
	return nil
}
