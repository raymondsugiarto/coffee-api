package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func FindAllAdmin(service admin.Service) fiber.Handler {
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

func AdminGetMyProfile(service admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID

		result, err := service.FindByUserID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateAdminProfileImage(service admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID
		profileImage, err := c.FormFile("profile_image")
		if err != nil {
			return status.New(status.BadRequest, err)
		}
		var profileImageUrlPath string = ""
		if profileImage != nil {
			profileImageUrlPath = fmt.Sprintf("./storage/admin/profile/%s", profileImage.Filename)

			if err := c.SaveFile(profileImage, profileImageUrlPath); err != nil {
				return status.New(status.BadRequest, err)
			}
		}

		if err := service.UpdateProfileImage(c.Context(), id, profileImageUrlPath); err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"profileImageUrl": profileImageUrlPath,
		})
	}
}

func UpdateAdminName(service admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID
		dto := new(entity.AdminDto)
		if err := c.BodyParser(dto); err != nil {
			return status.New(status.BadRequest, err)
		}

		response, err := service.UpdateName(c.Context(), id, dto)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func CreateAdminCompany(service admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dto := new(entity.CreateAdminCompany)
		if err := c.BodyParser(dto); err != nil {
			return status.New(status.BadRequest, err)
		}

		response, err := service.CreateAdminCompany(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func FindAllAdminByCompanyID(service admin.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("id")
		req := new(entity.FindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		result, err := service.FindAllByCompanyID(c.Context(), companyID, req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
