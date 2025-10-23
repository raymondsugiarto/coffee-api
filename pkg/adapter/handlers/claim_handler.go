package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/claim"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateClaim(service claim.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ClaimInputDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}

		dto := req.ToDto()
		claimAttachment(c, dto)

		createdClaim, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(createdClaim)
	}
}

func FindAllClaimByCustomer(service claim.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ClaimFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}

		userCredential := shared.GetUserCredential(c.Context())
		req.CustomerID = userCredential.CustomerID

		result, err := service.FindAll(c.Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindClaimByCompany(service claim.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ClaimFindAllRequest)
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

func FindClaimByID(service claim.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func claimAttachment(c *fiber.Ctx, dto *entity.ClaimDto) error {
	files := map[string]*string{
		"participantCard":         &dto.ParticipantCard,
		"identityCardFile":        &dto.IdentityCardFile,
		"taxIdentityFile":         &dto.TaxIdentityFile,
		"familyCardFile":          &dto.FamilyCardFile,
		"deathCertificateFile":    &dto.DeathCertificateFile,
		"guardianCertificateFile": &dto.GuardianCertificateFile,
		"medicalCertificateFile":  &dto.MedicalCertificateFile,
		"workCertificateFile":     &dto.WorkCertificateFile,
	}

	for fieldName, fieldPath := range files {
		file, err := c.FormFile(fieldName)
		if err != nil && err != fiber.ErrUnprocessableEntity {
			return status.New(status.BadRequest, fmt.Errorf("error retrieving file %s: %w", fieldName, err))
		}
		if file != nil {
			filePath := fmt.Sprintf("./storage/%s", file.Filename)
			if saveErr := c.SaveFile(file, filePath); saveErr != nil {
				return status.New(status.BadRequest, fmt.Errorf("error saving file %s: %w", fieldName, saveErr))
			}
			*fieldPath = filePath
		}
	}
	return nil
}
