package migrate

import (
	"fmt"
	"log"
	"os"
	"strconv"

	config "github.com/raymondsugiarto/coffee-api/config"
	"github.com/raymondsugiarto/coffee-api/pkg/infrastructure/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(schema string) *migrate.Migrate {
	cfg := config.GetConfig().Database.Main
	fmt.Printf("DB Config: %v\n", cfg)
	dbname := cfg.Dbname

	driver, err := database.GetDatabaseDriverMigration(cfg, schema)
	if err != nil {
		log.Fatalln("Failed to loading certificates: ", err)
		os.Exit(1)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/"+cfg.Adapter,
		dbname,
		driver,
	)
	if err != nil {
		log.Fatalln("Failed to loading certificates: ", err)
		os.Exit(1)
	}
	return m
}

// Migration :
func Migration(args []string, schema string) {
	fmt.Println("Run DB Migrate")
	m := Connect(schema)
	if len(args) >= 3 {
		if args[2] == "" {
			m.Steps(1)
		} else {
			step, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatalln("Args step must be number ", err)
				os.Exit(1)
			}
			if args[1] == "down" {
				step = step * -1
			}
			m.Steps(step)
		}
	} else {
		if len(args) >= 2 {
			if args[1] == "down" {
				migrationDown(m)
			} else {
				migrationUp(m)
			}
		} else {
			migrationUp(m)
		}
	}
}

func MigrateUpAll() {
	schema := config.GetConfig().Database.Main.Schema
	m := Connect(schema)
	fmt.Println("migrate up")
	err := m.Up()
	if err != nil {
		fmt.Printf("Error migration %v", err)
	}
}

func migrationUp(m *migrate.Migrate) {
	fmt.Println("migrate up")
	err := m.Up()
	if err != nil {
		fmt.Printf("Error migration %v", err)
	}
}

func migrationDown(m *migrate.Migrate) {
	fmt.Println("migrate down")
	err := m.Down()
	if err != nil {
		fmt.Printf("Error migration %v", err)
	}
}
