package customer

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/model"
	"github.com/raymondsugiarto/coffee-api/pkg/module/customer"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func UploadCustomer(service customer.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("document")
		if err != nil {
			return status.New(status.BadRequest, fmt.Errorf("file is required"))
		}

		// Buka file
		f, err := file.Open()
		if err != nil {
			return status.New(status.BadRequest, fmt.Errorf("cannot open uploaded file"))
		}
		defer f.Close()

		// Parsing Excel ke []*CustomerInputDto
		importExcel, err := service.ParseExcelToCustomers(c.Context(), f)
		if err != nil {
			return status.New(status.BadRequest, err)
		}

		if len(importExcel.Errors) > 0 {
			return status.New(status.BadRequest, importExcel.Errors...)
		}

		// Mapping ke []*CustomerDto
		var customersDto []*entity.CustomerDto
		for _, input := range importExcel.Data {
			dto := input.ToDto()
			dto.ApprovalStatus = model.ApprovalStatusApproved
			dto.SIMStatus = model.SIMStatusActive
			customersDto = append(customersDto, dto)
		}

		// Panggil service untuk create batch
		created, err := service.CreateBatch(c.Context(), customersDto)
		if err != nil {
			return status.New(status.InternalServerError, err)
		}

		return c.JSON(created)
	}
}
