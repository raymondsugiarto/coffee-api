package server

import (
	"os"

	config "github.com/raymondsugiarto/coffee-api/config"

	"log"
	"strconv"

	"github.com/raymondsugiarto/coffee-api/pkg/adapter/routes"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/database"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Rest struct {
}

func NewRest() *Rest {
	return &Rest{}
}

func (s *Rest) Initialize() {
	os.Setenv("TZ", "Asia/Jakarta")

	response.SetAppCode("D")
	cfg := config.GetConfig()

	app := fiber.New(fiber.Config{
		AppName:           cfg.Server.Rest.Name,
		ErrorHandler:      middleware.DefaultErrorHandler(),
		BodyLimit:         10 * 1024 * 1024,
		EnablePrintRoutes: false,
	})
	app.Use(cors.New())

	middleware.SetupValidator()
	initDatabase()

	routes.InitRouter(app)

	err := app.Listen(":" + strconv.Itoa(cfg.Server.Rest.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func initDatabase() {
	cfg := config.GetConfig()
	sqlConn, err := database.NewSQLConnection(cfg.Database.Main, cfg.Database.Main.Schema)
	if err != nil {
		log.Fatal(err)
	}
	database.DBConn = sqlConn.GetConn()
}
