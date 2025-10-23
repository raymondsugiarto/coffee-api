package handlers

import (
	"github.com/gofiber/fiber/v2"
	entity "github.com/raymondsugiarto/coffee-api/pkg/entity/investment"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	investmentdistribution "github.com/raymondsugiarto/coffee-api/pkg/module/investment/investment_distribution"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

func CreateBatchInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentDistributionInputBatchDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := make([]*entity.InvestmentDistributionDto, len(request.Data))
		for i, data := range request.Data {
			dto[i] = data.ToDto()
		}

		result, err := service.CreateBatch(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func CreateInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentDistributionInputDto)

		if err := c.BodyParser(request); err != nil {
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

func GetAllInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.InvestmentFindAllRequest)
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

func GetInvestmentDistributionByID(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		result, err := service.FindByID(c.Context(), id)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetInvestmentDistributionByCompany(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		companyID := c.Params("companyId")

		result, err := service.FindByCompanyID(c.Context(), companyID)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func GetInvestmentDistributionByParticipant(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		participantID := c.Params("participantId")

		result, err := service.FindByParticipantID(c.Context(), participantID)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}

func UpdateInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		request := new(entity.InvestmentDistributionInputDto)

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

func DeleteInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := service.Delete(c.Context(), id); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func CreateOrUpdateInvestmentDistribution(service investmentdistribution.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(entity.InvestmentDistributionInputBatchDto)

		if err := c.BodyParser(request); err != nil {
			return status.New(status.BadRequest, err)
		}

		if err := middleware.AppValidator.Validate(request); err != nil {
			return err
		}

		dto := make([]*entity.InvestmentDistributionDto, len(request.Data))
		for i, data := range request.Data {
			dto[i] = data.ToDto()
		}

		result, err := service.CreateOrUpdate(c.Context(), dto)
		if err != nil {
			return err
		}

		return c.JSON(result)
	}
}
