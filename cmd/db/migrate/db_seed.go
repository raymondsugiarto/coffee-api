package migrate

// Seed :
// func RunSeeding(args []string, schema string, seedName string) {
// 	fmt.Printf("Start seeding %s on %s \n", seedName, schema)
// 	cfg := config.GetConfig()
// 	logConfig := logger.Configuration{
// 		EnableConsole:     cfg.Logger.Console.Enable,
// 		ConsoleJSONFormat: cfg.Logger.Console.JSON,
// 		ConsoleLevel:      cfg.Logger.Console.Level,
// 		EnableFile:        cfg.Logger.File.Enable,
// 		FileJSONFormat:    cfg.Logger.File.JSON,
// 		FileLevel:         cfg.Logger.File.Level,
// 		FileLocation:      cfg.Logger.File.Path,
// 	}

// 	if err := logger.NewLogger(logConfig, logger.InstanceZapLogger); err != nil {
// 		log.Fatalf("Could not instantiate log %v", err)
// 	}

// 	driver, err := db.NewSQLConnection(cfg.Database.Main, schema)
// 	if err != nil {
// 		logger.Fatalf("Failed to connect DB %v", err)
// 	}

// 	seeds := make([]seed.Seed, 0)

// 	switch seedName {
// 	case "admin" :
// 		seeds = append(seeds, seed.RunAdminSeeds()...)
// 		seeds = append(seeds, seed.RunRoleSeeds()...)
// 		seeds = append(seeds, seed.RunUserRolesSeeds()...)
// 		seeds = append(seeds, seed.RunMenus()...)
// 		seeds = append(seeds, seed.RunRolePermissions()...)
// 	case "location" :
// 		seeds = append(seeds, seed.RunLocationSeeds()...)
// 	case "countries" :
// 		seeds = append(seeds, seed.RunCountriesSeed()...)
// 	case "work" :
// 		seeds = append(seeds, seed.RunWorkSeeds()...)
// 	case "all" :
// 		seeds = append(seeds, seed.RunAdminSeeds()...)
// 		seeds = append(seeds, seed.RunRoleSeeds()...)
// 		seeds = append(seeds, seed.RunUserRolesSeeds()...)
// 		seeds = append(seeds, seed.RunMenus()...)
// 		seeds = append(seeds, seed.RunRolePermissions()...)
// 		seeds = append(seeds, seed.RunLocationSeeds()...)
// 		seeds = append(seeds, seed.RunCountriesSeed()...)
// 		seeds = append(seeds, seed.RunWorkSeeds()...)
// 	default :
// 		fmt.Println("Please specify a seed name")
// 	}

// 	for _, seed := range seeds {
// 		if err := seed.Run(driver.GetConn()); err != nil {
// 			logger.Fatalf("Running seed '%s', failed with error: %s", seedName, err)
// 		}
// 	}

// 	fmt.Printf("Seeding completed. %v seeds running\n", len(seeds))

// }
