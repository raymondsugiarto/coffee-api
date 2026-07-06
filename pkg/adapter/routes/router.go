package routes

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	handlers "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers"
	ha "github.com/raymondsugiarto/coffee-api/pkg/adapter/handlers/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/database"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware/organization"
	"github.com/raymondsugiarto/coffee-api/pkg/module/admin"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication"
	"github.com/raymondsugiarto/coffee-api/pkg/module/authentication/token"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/driver"
	"github.com/raymondsugiarto/coffee-api/pkg/module/item"
	"github.com/raymondsugiarto/coffee-api/pkg/module/order"
	orderitem "github.com/raymondsugiarto/coffee-api/pkg/module/order/order_item"
	stocksession "github.com/raymondsugiarto/coffee-api/pkg/module/stock_session"
	"github.com/raymondsugiarto/coffee-api/pkg/module/user"
	usercredential "github.com/raymondsugiarto/coffee-api/pkg/module/user-credential"

	"github.com/gofiber/fiber/v2"
)

func InitRouter(app fiber.Router) {
	dbConn := database.DBConn
	app.Use(logger.New(), organization.New(), middleware.DefaultResponseHandler())

	// User Credential
	userCredentialRepo := usercredential.NewRepository(dbConn)
	userCredentialService := usercredential.NewService(userCredentialRepo)

	// User
	userRepo := user.NewRepository(dbConn)
	userService := user.NewService(userRepo, userCredentialService)

	// Admin
	adminRepo := admin.NewRepository(dbConn)
	adminService := admin.NewService(adminRepo)

	tokenService := token.NewService()

	authenticationService := authentication.NewService(
		userCredentialService, tokenService, adminService,
	)

	// Item
	companyRepo := company.NewRepository(dbConn)
	companyService := company.NewService(companyRepo)

	// Item
	itemRepo := item.NewRepository(dbConn)
	itemService := item.NewService(itemRepo, companyService)

	// Order
	orderRepo := order.NewRepository(dbConn)
	orderService := order.NewService(orderRepo, companyService)

	// Order
	orderItemRepo := orderitem.NewRepository(dbConn)
	orderItemService := orderitem.NewService(orderItemRepo, companyService)

	// Driver (employees filtered)
	driverService := driver.NewService(dbConn)

	// Stock Session (with embedded item picker service)
	stockSessionRepo := stocksession.NewRepository(dbConn)
	stockSessionService := stocksession.NewService(stockSessionRepo, dbConn)
	stockSessionItemService := stocksession.NewItemService(dbConn)

	// Middleware
	// api := app.Group("/api", middleware.Protected())
	auth := app.Group("/api/auth")
	AuthRouter(auth, userService, authenticationService)

	api := app.Group("/api/", middleware.Protected())
	ItemRouter(api, itemService)
	OrderRouter(api, orderService)
	OrderItemRouter(api, orderItemService)
	ProductRouter(api, stockSessionItemService)
	DriverRouter(api, driverService)
	StockSessionRouter(api, stockSessionService, stockSessionItemService)
}

func AuthRouter(app fiber.Router,
	service user.Service,
	authService authentication.Service,
) {
	app.Post("/sign-in", ha.SignIn(authService))
}

func ItemRouter(app fiber.Router,
	itemService item.Service,
) {
	app.Get("/items", handlers.FindAllItems(itemService))
}

func OrderRouter(app fiber.Router,
	orderService order.Service,
) {
	app.Post("/orders", handlers.CreateOrder(orderService))
	app.Get("/orders", handlers.FindAllMyOrders(orderService))
	app.Get("/orders/count", handlers.CountMyOrders(orderService))
}

func OrderItemRouter(app fiber.Router,
	orderItemService orderitem.Service,
) {
	app.Get("/order-items/count", handlers.CountMyOrderItems(orderItemService))
}

func ProductRouter(app fiber.Router, itemService stocksession.ItemService) {
	app.Get("/products", handlers.FindAllStockSessionItems(itemService))
	app.Get("/products/:id", handlers.GetStockSessionItem(itemService))
}

func DriverRouter(app fiber.Router, driverService driver.Service) {
	app.Get("/employees", handlers.FindAllDrivers(driverService))
}

func StockSessionRouter(app fiber.Router, ssService stocksession.Service, itemService stocksession.ItemService) {
	app.Post("/stock-session/open", handlers.OpenStockSession(ssService))
	app.Get("/stock-session", handlers.FindAllStockSessions(ssService))
	app.Get("/stock-session/today", handlers.GetTodayStockSession(ssService))
	app.Get("/stock-session/:id", handlers.GetStockSession(ssService))
	app.Put("/stock-session/:id", handlers.UpdateStockSession(ssService))
	app.Post("/stock-session/:id/close", handlers.CloseStockSession(ssService))

	// Item picker (reuses existing `item` table)
	app.Get("/products", handlers.FindAllStockSessionItems(itemService))
	app.Get("/products/:id", handlers.GetStockSessionItem(itemService))

	// Reports
	app.Get("/report/dashboard", handlers.GetDashboard(ssService))
	app.Get("/report/daily", handlers.GetDailyReport(ssService))
	app.Get("/report/monthly", handlers.GetMonthlyReport(ssService))
	app.Get("/report/top-products", handlers.GetTopProducts(ssService))
	app.Get("/report/employee-performance", handlers.GetEmployeePerformance(ssService))
}
