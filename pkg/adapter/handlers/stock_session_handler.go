package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/entity"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	stocksession "github.com/raymondsugiarto/coffee-api/pkg/module/stock_session"
	shared "github.com/raymondsugiarto/coffee-api/pkg/shared/context"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

// ============ Item picker (reuses existing `item` table) ============

func FindAllStockSessionItems(service stocksession.ItemService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.ItemFindAllRequest)
		if err := c.QueryParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		result, err := service.FindAllItems(c.Context(), req)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func GetStockSessionItem(service stocksession.ItemService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.GetItem(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// GetStockSessionItemChildren returns items whose parent_id matches any of the
// supplied parentIds. Supports comma-separated `parentIds` query param plus an
// `includeInactive` flag. Used by the close-session UI to auto-expand parents.
func GetStockSessionItemChildren(service stocksession.ItemService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(c.Query("parentIds"))
		var ids []string
		for _, p := range strings.Split(raw, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				ids = append(ids, p)
			}
		}
		includeInactive := c.Query("includeInactive") == "true"
		result, err := service.GetItemChildren(c.Context(), ids, includeInactive)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

type setItemParentRequest struct {
	ParentID string   `json:"parentId" validate:"required"`
	ChildIDs []string `json:"childIds" validate:"required,min=1,dive,required"`
}

// SetStockSessionItemParent bulk-updates the parent_id of every child in
// childIds to point at parentId. The relation is keyed on UUID, not on code.
// Used by the catalog admin UI to declare which items roll up under which
// parent so the close-session picker can auto-expand variants.
func SetStockSessionItemParent(service stocksession.ItemService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(setItemParentRequest)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		if err := service.SetItemParent(c.Context(), req.ParentID, req.ChildIDs); err != nil {
			return status.New(status.BadRequest, err)
		}
		return c.JSON(fiber.Map{
			"updated":  len(req.ChildIDs),
			"parentId": req.ParentID,
		})
	}
}

// ============ Open ============

func OpenStockSession(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.OpenStockSessionInputDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}

		dto := &entity.StockSessionDto{
			EmployeeID: req.EmployeeID,
			Date:       req.Date,
			Notes:      req.Notes,
		}
		for _, it := range req.Items {
			internal := it.ToStockSessionItemInputDto()
			dto.Items = append(dto.Items, *internal.ToDto())
		}

		userCred := shared.GetUserCredential(c.Context())
		actorID := ""
		if userCred != nil {
			actorID = userCred.AdminID
		}

		result, err := service.Open(c.Context(), dto, actorID)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

// ============ Get / FindAll ============

func GetStockSession(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := service.Get(c.Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func FindAllStockSessions(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := new(entity.StockSessionFindAllRequest)
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

func GetTodayStockSession(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		employeeID := c.Query("employeeId")
		date := c.Query("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		if employeeID == "" {
			return status.New(status.BadRequest, errors.New("employeeId is required"))
		}
		result, err := service.GetByEmployeeDate(c.Context(), employeeID, date)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// ============ Update (still open) ============

func UpdateStockSession(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.OpenStockSessionInputDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		dto := &entity.StockSessionDto{
			ID:         id,
			EmployeeID: req.EmployeeID,
			Date:       req.Date,
			Notes:      req.Notes,
		}
		for _, it := range req.Items {
			internal := it.ToStockSessionItemInputDto()
			dto.Items = append(dto.Items, *internal.ToDto())
		}

		userCred := shared.GetUserCredential(c.Context())
		actorID := ""
		if userCred != nil {
			actorID = userCred.AdminID
		}

		result, err := service.Update(c.Context(), dto, actorID)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// ============ Close ============

func CloseStockSession(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		req := new(entity.CloseStockSessionInputDto)
		if err := c.BodyParser(req); err != nil {
			return status.New(status.BadRequest, err)
		}
		if err := middleware.AppValidator.Validate(req); err != nil {
			return err
		}
		dto := &entity.StockSessionDto{
			Notes: req.Notes,
		}
		for _, it := range req.Items {
			internal := it.ToStockSessionItemInputDto()
			dto.Items = append(dto.Items, *internal.ToDto())
		}
		for _, p := range req.Payments {
			dto.Payments = append(dto.Payments, *p.ToDto())
		}
		for _, a := range req.Adjustments {
			dto.Adjustments = append(dto.Adjustments, *a.ToDto())
		}

		userCred := shared.GetUserCredential(c.Context())
		actorID := ""
		if userCred != nil {
			actorID = userCred.AdminID
		}

		result, err := service.Close(c.Context(), id, dto, actorID)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// ============ Reports / Dashboard ============

func GetDashboard(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		result, err := service.GetDashboard(c.Context())
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func GetDailyReport(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		date := c.Query("date")
		result, err := service.GetDailyReport(c.Context(), date)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func GetMonthlyReport(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		year, _ := strconv.Atoi(c.Query("year"))
		month, _ := strconv.Atoi(c.Query("month"))
		result, err := service.GetMonthlyReport(c.Context(), year, month)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func GetTopProducts(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		limit, _ := strconv.Atoi(c.Query("limit"))
		result, err := service.GetTopProducts(c.Context(), from, to, limit)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

func GetEmployeePerformance(service stocksession.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		from := c.Query("from")
		to := c.Query("to")
		result, err := service.GetEmployeePerformance(c.Context(), from, to)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}
