package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/utils"
)

func SignUpCustomer(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.SignUpCustomerInputDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()

		response, err := service.SignUpCustomer(c.Context(), dto)
		if err != nil {
			log.Errorf("errorSignUpCustomer: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func SignUpCompany(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.SignUpCompanyInputDto)
		if err := c.BodyParser(request); err != nil {
			log.Errorf("errorSignUpCompany: %v", err)
			return status.New(status.BadRequest, err)
		}

		companyCode, err := gonanoid.Generate(utils.ALPHA_NUMERIC, 5)
		if err != nil {
			log.Errorf("errorGenerateCompanyCode: %v", err)
			return status.New(status.BadRequest, err)
		}

		dto := request.ToDto()
		dto.CompanyCode = companyCode
		companyAttachment(c, dto)

		response, err := service.SignUpCompany(c.Context(), dto)
		if err != nil {
			log.Errorf("errorSignUpCompany: %v", err)
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func SignIn(service authentication.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.LoginRequestDto)
		if err := c.BodyParser(&request); err != nil {
			return status.New(status.BadRequest, err)
		}

		response, err := service.SignIn(c.Context(), request)
		if err != nil {
			fmt.Printf("errorSignIn: %v\n", err)
			return err
		}

		return c.JSON(response)
	}
}

func companyAttachment(c *fiber.Ctx, dto *entity.SignUpCompanyDto) error {
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
