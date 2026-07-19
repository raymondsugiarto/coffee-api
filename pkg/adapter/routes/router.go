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
	cashdebt "github.com/raymondsugiarto/coffee-api/pkg/module/cash_debt"
	"github.com/raymondsugiarto/coffee-api/pkg/module/company"
	"github.com/raymondsugiarto/coffee-api/pkg/module/driver"
	"github.com/raymondsugiarto/coffee-api/pkg/module/item"
	itemcategory "github.com/raymondsugiarto/coffee-api/pkg/module/item_category"
	"github.com/raymondsugiarto/coffee-api/pkg/module/order"
	orderitem "github.com/raymondsugiarto/coffee-api/pkg/module/order/order_item"
	"github.com/raymondsugiarto/coffee-api/pkg/module/payroll"
	salarycomponent "github.com/raymondsugiarto/coffee-api/pkg/module/salary_component"
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
	itemService := item.NewService(itemRepo)

	// Item Category
	itemCategoryRepo := itemcategory.NewRepository(dbConn)
	itemCategoryService := itemcategory.NewService(itemCategoryRepo)

	// Salary Component
	salaryComponentRepo := salarycomponent.NewRepository(dbConn)
	salaryComponentService := salarycomponent.NewService(salaryComponentRepo)

	// Payroll (employee_salary + employee_salary_component)
	payrollRepo := payroll.NewRepository(dbConn)
	payrollService := payroll.NewService(payrollRepo)

	// Cash Debt (driver cash advances ledger)
	cashDebtRepo := cashdebt.NewRepository(dbConn)
	cashDebtService := cashdebt.NewService(cashDebtRepo)

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
	ItemCategoryRouter(api, itemCategoryService)
	SalaryComponentRouter(api, salaryComponentService)
	PayrollRouter(api, payrollService)
	CashDebtRouter(api, cashDebtService)
	CompanyRouter(api, companyService)
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
	app.Get("/items/:id", handlers.FindOneItem(itemService))
	app.Post("/items", handlers.CreateItem(itemService))
	app.Put("/items/:id", handlers.UpdateItem(itemService))
	app.Delete("/items/:id", handlers.DeleteItem(itemService))
}

func ItemCategoryRouter(app fiber.Router,
	itemCategoryService itemcategory.Service,
) {
	app.Get("/item-categories", handlers.FindAllItemCategories(itemCategoryService))
	app.Get("/item-categories/:id", handlers.FindOneItemCategory(itemCategoryService))
	app.Post("/item-categories", handlers.CreateItemCategory(itemCategoryService))
	app.Put("/item-categories/:id", handlers.UpdateItemCategory(itemCategoryService))
	app.Delete("/item-categories/:id", handlers.DeleteItemCategory(itemCategoryService))
}

func SalaryComponentRouter(app fiber.Router,
	salaryComponentService salarycomponent.Service,
) {
	app.Get("/salary-components", handlers.FindAllSalaryComponents(salaryComponentService))
	app.Get("/salary-components/:id", handlers.FindOneSalaryComponent(salaryComponentService))
	app.Post("/salary-components", handlers.CreateSalaryComponent(salaryComponentService))
	app.Put("/salary-components/:id", handlers.UpdateSalaryComponent(salaryComponentService))
	app.Delete("/salary-components/:id", handlers.DeleteSalaryComponent(salaryComponentService))
}

// CompanyRouter exposes the read-only company list used by the
// SelectCompany dropdown (and any future admin picker).
// Write operations are intentionally NOT wired here: the legacy
// company handler in pkg/adapter/handlers/company uses a multipart
// attachment flow that lives outside this scope.
func CompanyRouter(app fiber.Router,
	companyService company.Service,
) {
	app.Get("/companies", handlers.FindAllCompanies(companyService))
}

// PayrollRouter wires the payroll run lifecycle:
//
//	POST /payroll/simulate — read-only preview
//	POST /payroll          — persist the approved run
//	GET  /payroll          — list saved runs
//	GET  /payroll/:id      — one run with components
func PayrollRouter(app fiber.Router,
	payrollService payroll.Service,
) {
	app.Post("/payroll/simulate", handlers.SimulatePayroll(payrollService))
	app.Post("/payroll", handlers.SavePayroll(payrollService))
	app.Get("/payroll", handlers.FindAllPayrolls(payrollService))
	app.Get("/payroll/:id", handlers.FindOnePayroll(payrollService))
}

// CashDebtRouter wires the driver cash-advance ledger CRUD.
func CashDebtRouter(app fiber.Router,
	cashDebtService cashdebt.Service,
) {
	app.Get("/cash-debts", handlers.FindAllCashDebts(cashDebtService))
	app.Get("/cash-debts/:id", handlers.FindOneCashDebt(cashDebtService))
	app.Post("/cash-debts", handlers.CreateCashDebt(cashDebtService))
	app.Put("/cash-debts/:id", handlers.UpdateCashDebt(cashDebtService))
	app.Delete("/cash-debts/:id", handlers.DeleteCashDebt(cashDebtService))
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
	app.Get("/products/children", handlers.GetStockSessionItemChildren(itemService))
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
	app.Get("/products/children", handlers.GetStockSessionItemChildren(itemService))
	app.Get("/products/:id", handlers.GetStockSessionItem(itemService))
	app.Post("/products/parent", handlers.SetStockSessionItemParent(itemService))

	// Reports
	app.Get("/report/dashboard", handlers.GetDashboard(ssService))
	app.Get("/report/daily", handlers.GetDailyReport(ssService))
	app.Get("/report/monthly", handlers.GetMonthlyReport(ssService))
	app.Get("/report/top-products", handlers.GetTopProducts(ssService))
	app.Get("/report/employee-performance", handlers.GetEmployeePerformance(ssService))
}
