/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/akkaraju-satvik/dbmap/config"
	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var createMigrationErrorStrings = map[string]string{
	"config":         "Error reading config file\n %s",
	"migrationDir":   "Error creating migration directory\n %s",
	"migrationFiles": "Error creating migration files\n %s",
	"database":       "Error connecting to the database\n %s",
	"migration":      "Error creating migration\n %s",
}

var createMigrationCmd = &cobra.Command{
	Use:   "create-migration",
	Short: "Create a new migration in the migrations directory specified in the config file",
	Long: `Creates a new migration in the migrations directory specified in the config file.

The migration will have an up.sql and down.sql file with placeholders for the migration queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config-file").Value.String()
		config, err := config.Read(configFilePath)
		if err != nil {
			color.Red(createMigrationErrorStrings["config"], err.Error())
			os.Exit(1)
		}
		err = createMigration(config)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	},
}

func createMigration(config config.Config) error {
	currentTime := fmt.Sprintf("%d", time.Now().Unix())
	migrationDir := config.MigrationsDir + "/" + currentTime
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf(createMigrationErrorStrings["migrationDir"], err.Error())
	}
	err := createMigrationFiles(migrationDir)
	if err != nil {
		return fmt.Errorf(createMigrationErrorStrings["migrationFiles"], err.Error())
	}
	db, err := sql.Open("postgres", config.DbURL)
	if err != nil {
		removeMigrationSetup(migrationDir)
		return fmt.Errorf(createMigrationErrorStrings["database"], err.Error())
	}
	_, err = db.Exec(queries.CreateMigrationQuery, currentTime)
	if err != nil {
		removeMigrationSetup(migrationDir)
		return fmt.Errorf(createMigrationErrorStrings["migration"], err.Error())
	}
	db.Close()
	color.Green("Migration created successfully")
	return nil
}

func init() {
	createMigrationCmd.Flags().StringP("config-file", "c", "dbmap.config.yaml", "config file to use")
	rootCmd.AddCommand(createMigrationCmd)
}

func createMigrationFiles(migrationDir string) error {
	upFile, err := os.Create(migrationDir + "/up.sql")
	if err != nil {
		removeMigrationSetup(migrationDir)
		return fmt.Errorf("error creating migration file\n %s", err.Error())
	}
	fmt.Fprintf(upFile, "-- Write your UP migration here\n\n")
	upFile.Close()
	downFile, err := os.Create(migrationDir + "/down.sql")
	if err != nil {
		removeMigrationSetup(migrationDir)
		return fmt.Errorf("error creating migration file\n %s", err.Error())
	}
	fmt.Fprintf(downFile, "-- Write your DOWN migration here\n\n")
	downFile.Close()
	return nil
}

func removeMigrationSetup(migrationDir string) {
	os.RemoveAll(migrationDir)
}
