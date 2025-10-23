package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// func FindAllAdmin(service company.Service) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		req := new(entity.FindAllRequest)
// 		if err := c.QueryParser(req); err != nil {
// 			return status.New(status.BadRequest, err)
// 		}

// 		result, err := service.FindAll(c.Context(), req)
// 		if err != nil {
// 			return err
// 		}

// 		return c.JSON(result)
// 	}
// }

func CompanyGetMyProfile(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := shared.GetUserCredential(c.Context()).UserID

		result, err := service.FindByUserID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func CreateCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.CompanyInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := request.ToDto()
		companyAttachment(c, dto)

		result, err := service.Create(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindAllCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.FindAllRequest)
		if err := c.QueryParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		result, err := service.FindAll(c.Context(), request)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func FindCompanyByID(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func DeleteCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c.Context(), id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.CompanyInputDto)
		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.ID = id
		companyAttachment(c, dto)

		result, err := service.Update(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func companyAttachment(c *fiber.Ctx, dto *entity.CompanyDto) error {
	// companyCode := dto.CompanyCode
	cooperationAgreement, errCooperationAgreement := c.FormFile("cooperationAgreement")
	aktaPerusahaan, errAktaPerusahaan := c.FormFile("aktaPerusahaan")
	nibFile, errNibFile := c.FormFile("nibFile")
	tdp, errTdp := c.FormFile("tdp")
	ktp, errKtp := c.FormFile("ktp")
	npwpPerusahaan, errNpwpPerusahaan := c.FormFile("npwpPerusahaan")
	suratKuasa, errSuratKuasa := c.FormFile("suratKuasa")

	if cooperationAgreement != nil {
		cooperationAgreementPath := fmt.Sprintf("./storage/%s", cooperationAgreement.Filename)
		// Save the files to disk:
		if err := c.SaveFile(cooperationAgreement, cooperationAgreementPath); err != nil {
			return status.New(status.BadRequest, errCooperationAgreement)
		}
		dto.CooperationAgreement = cooperationAgreementPath
	}
	if aktaPerusahaan != nil {
		aktaPerusahaanPath := fmt.Sprintf("./storage/%s", aktaPerusahaan.Filename)
		// Save the files to disk:
		if err := c.SaveFile(aktaPerusahaan, aktaPerusahaanPath); err != nil {
			return status.New(status.BadRequest, errAktaPerusahaan)
		}
		dto.AktaPerusahaan = aktaPerusahaanPath
	}
	if nibFile != nil {
		nibFilePath := fmt.Sprintf("./storage/%s", nibFile.Filename)
		// Save the files to disk:
		if err := c.SaveFile(nibFile, nibFilePath); err != nil {
			return status.New(status.BadRequest, errNibFile)
		}
		dto.NIBFile = nibFilePath
	}
	if tdp != nil {
		tdpPath := fmt.Sprintf("./storage/%s", tdp.Filename)
		// Save the files to disk:
		if err := c.SaveFile(tdp, tdpPath); err != nil {
			return status.New(status.BadRequest, errTdp)
		}
		dto.TDP = tdpPath
	}
	if ktp != nil {
		ktpPath := fmt.Sprintf("./storage/%s", ktp.Filename)
		// Save the files to disk:
		if err := c.SaveFile(ktp, ktpPath); err != nil {
			return status.New(status.BadRequest, errKtp)
		}
		dto.KTP = ktpPath
	}
	if npwpPerusahaan != nil {
		npwpPerusahaanPath := fmt.Sprintf("./storage/%s", npwpPerusahaan.Filename)
		// Save the files to disk:
		if err := c.SaveFile(npwpPerusahaan, npwpPerusahaanPath); err != nil {
			return status.New(status.BadRequest, errNpwpPerusahaan)
		}
		dto.NPWPPerusahaan = npwpPerusahaanPath
	}
	if suratKuasa != nil {
		suratKuasaPath := fmt.Sprintf("./storage/%s", suratKuasa.Filename)
		// Save the files to disk:
		if err := c.SaveFile(suratKuasa, suratKuasaPath); err != nil {
			return status.New(status.BadRequest, errSuratKuasa)
		}
		dto.SuratKuasa = suratKuasaPath
	}
	return nil
}
