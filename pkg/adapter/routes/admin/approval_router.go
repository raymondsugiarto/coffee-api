package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	"github.com/raymondsugiarto/coffee-api/pkg/module/approval"
)

func ApprovalRouter(app fiber.Router, approvalSvc approval.Service) {
	app.Get("", handlers.FindAllApproval(approvalSvc))
	app.Get("/:id", handlers.FindApprovalByID(approvalSvc))
	app.Post("/:id/confirmation", handlers.ConfirmationApproval(approvalSvc))
}
