package cmd

import (
	"os"

	"github.com/akkaraju-satvik/dbmap/config"
	"github.com/akkaraju-satvik/dbmap/migrations"
	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

type Migration struct {
	MigrationId     string `db:"migration_id"`
	MigrationName   string `db:"migration_name"`
	MigrationTime   string `db:"migration_time"`
	MigrationStatus string `db:"migration_status"`
}

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
		err = applyMigrations(configuration)
		if err != nil {
			color.Red("Error applying migrations\n %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	applyMigrationsCmd.Flags().StringP("config-file", "c", "dbmap.config.yaml", "Path to the config file")
	applyMigrationsCmd.Flags().StringP("migration-type", "t", "up", "Type of migration to apply")
	rootCmd.AddCommand(applyMigrationsCmd)
}

func applyMigrations(config config.Config) error {
	conn, err := sqlx.Open("postgres", config.DbURL)
	if err != nil {
		return err
	}

	var migrationList []Migration
	err = conn.Select(&migrationList, queries.GetMigrations)
	if err != nil {
		return err
	}
	for _, migration := range migrationList {
		if migration.MigrationStatus == "APPLIED" {
			continue
		}
		migrationQuery, err := migrations.GetQuery(config.MigrationsDir, migration.MigrationName, "up")
		if err != nil {
			return err
		}
		color.Yellow("Applying migration %s", migration.MigrationName)
		_, err = conn.Exec(migrationQuery)
		if err != nil {
			return err
		}
		_, err = conn.Exec(queries.UpdateMigrationStatus, migration.MigrationId)
		if err != nil {
			return err
		}
		_, err = conn.Exec(queries.InsertMigrationQuery, migration.MigrationId, migrationQuery)
		if err != nil {
			return err
		}
		color.Green("Migration %s applied successfully", migration.MigrationName)
	}
	return nil
}
