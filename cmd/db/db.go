package db

import (
	"fmt"

	"github.com/raymondsugiarto/coffee-api/cmd/db/migrate"

	"github.com/spf13/cobra"
)

// DBCmd represents the db command
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		schema, _ := cmd.Flags().GetString("schema")
		// seedName, _ := cmd.Flags().GetString("name")
		fmt.Println(args)
		if args[0] == "migrate" {
			//fmt.Println(schema)
			migrate.Migration(args, schema)
		} else if args[0] == "seed" {
			// migrate.RunSeeding(args, schema, seedName)
		}
	},
}
