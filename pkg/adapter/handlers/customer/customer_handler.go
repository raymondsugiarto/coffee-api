package customer

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.CustomerInputDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		if c.Get("Accept") == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
			excelBytes, err := service.GenerateExcel(c.Context(), req, result)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Set("Content-Disposition", "attachment; filename=customers.xlsx")
			return c.SendStream(bytes.NewReader(excelBytes))
		}

		return c.JSON(result)
	}
}

func FindAllMyReferral(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		req.UserID = shared.GetUserCredential(c.Context()).UserID

		result, err := service.FindAllMyReferral(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllByCompany(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.CustomerFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAllByCompany(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindCustomerByID(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByIDWithScope(c.Context(), id, []string{"complete"})
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
func DeleteCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.CustomerInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id
		customerAttachment(c, dto)

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func SuspendStatusCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.CustomerInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
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

func ChangePassword(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.PasswordChangeInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		request.UserID = shared.GetUserCredential(c.Context()).UserID

		if err := middleware.AppValidator.Validate(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := service.CustomerChangePassword(c.Context(), request); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)

	}
}

func customerAttachment(c *fiber.Ctx, dto *entity.CustomerDto) error {
	identityCardFile, errIdentityCardFile := c.FormFile("identityCardFile")
	customerPhoto, errCustomerPhoto := c.FormFile("customerPhoto")
	taxIdentityCardFile, errTaxIdentityCardFile := c.FormFile("taxIdentityCardFile")

	if identityCardFile != nil {
		identityCardFilePath := fmt.Sprintf("./storage/%s", identityCardFile.Filename)

		if err := c.SaveFile(identityCardFile, identityCardFilePath); err != nil {
			return status.New(status.BadRequest, errIdentityCardFile)
		}
		dto.IdentityCardFile = identityCardFilePath
	}
	if customerPhoto != nil {
		customerPhotoPath := fmt.Sprintf("./storage/%s", customerPhoto.Filename)

		if err := c.SaveFile(customerPhoto, customerPhotoPath); err != nil {
			return status.New(status.BadRequest, errCustomerPhoto)
		}
		dto.CustomerPhoto = customerPhotoPath
	}
	if taxIdentityCardFile != nil {
		taxIdentityCardFilePath := fmt.Sprintf("./storage/%s", taxIdentityCardFile.Filename)

		if err := c.SaveFile(taxIdentityCardFile, taxIdentityCardFilePath); err != nil {
			return status.New(status.BadRequest, errTaxIdentityCardFile)
		}
		dto.TaxIdentityCardFile = taxIdentityCardFilePath
	}
	return nil
}
