package cmd

import (
	"fmt"
	"os"

	"github.com/akkaraju-satvik/dbmap/config"
	"github.com/akkaraju-satvik/dbmap/migrations"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var applyMigrationsCmd = &cobra.Command{
	Use:   "apply-migrations",
	Short: "Apply migrations in the migrations directory specified in the config file",
	Long: `Applies migrations in the migrations directory specified in the config file.

The migrations will be applied in the order of their creation.`,

	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config-file").Value.String()
		configuration, err := config.Read(configFilePath)
		if err != nil {
			color.Red("Error reading config file\n %s", err.Error())
			os.Exit(1)
		}
		migrationType := cmd.Flag("type").Value.String()
		if migrationType != "UP" && migrationType != "DOWN" {
			fmt.Println("Invalid values for --type flag:\nExpected Values:\n\tUP\n\tDOWN")
			os.Exit(1)
		}
		err = applyMigrations(configuration, migrationType)
		if err != nil {
			color.Red("Error applying migrations\n %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	applyMigrationsCmd.Flags().StringP("config-file", "c", "dbmap.config.yaml", "Path to the config file")
	applyMigrationsCmd.Flags().StringP("type", "t", "UP", "Type of migration to apply")
	rootCmd.AddCommand(applyMigrationsCmd)
}

func applyMigrations(config config.Config, migrationType string) error {

	var migrationList []migrations.Migration
	err := migrations.GetMigrations(config.DbURL, migrationType, &migrationList)
	if err != nil {
		return err
	}
	for _, migration := range migrationList {
		migrationQuery, err := migrations.GetQuery(config.MigrationsDir, migration.MigrationName, migrationType)
		if err != nil {
			return err
		}
		color.Yellow("Applying migration %s", migration.MigrationName)
		err = migrations.ApplyMigration(config.DbURL, migrationQuery, migrationType, migration)
		if err != nil {
			return err
		}
		color.Green("Migration %s applied successfully", migration.MigrationName)
	}
	return nil
}
