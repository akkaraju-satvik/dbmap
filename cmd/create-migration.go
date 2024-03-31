/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/akkaraju-satvik/dbmap/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var createMigrationCmd = &cobra.Command{
	Use:   "create-migration",
	Short: "Create a new migration in the migrations directory specified in the config file",
	Long: `Creates a new migration in the migrations directory specified in the config file.

The migration will have an up.sql and down.sql file with placeholders for the migration queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config-file").Value.String()
		config, err := utils.ReadConfig(configFilePath)
		if err != nil {
			color.Red("Error reading config file\n %s", err.Error())
			os.Exit(1)
		}
		utils.HandleEmptyValuesInConfig(config)
		currentTime := fmt.Sprintf("%d", time.Now().Unix())
		migrationDir := config.MigrationsDir + "/" + currentTime
		if err = os.MkdirAll(migrationDir, 0755); err != nil {
			color.Red("Error creating migration directory\n %s", err.Error())
			os.Exit(1)
		}
		err = createUpAndDownMigrationFiles(migrationDir)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		db, err := sql.Open("postgres", config.DBConnection)
		if err != nil {
			handleErrorWithRemoveMigration(migrationDir, fmt.Sprintf("Error connecting to the database\n %s", err.Error()))
		}
		_, err = db.Exec(queries.CreateMigrationQuery, currentTime)
		if err != nil {
			handleErrorWithRemoveMigration(migrationDir, fmt.Sprintf("Error creating migration\n %s", err.Error()))
		}
		db.Close()
		color.Green("Migration created successfully")
	},
}

func init() {
	createMigrationCmd.Flags().StringP("config-file", "c", "dbmap.json", "config file to use")
	rootCmd.AddCommand(createMigrationCmd)
}

func removeMigrationSetup(migrationDir string) {
	os.RemoveAll(migrationDir)
}

func createUpAndDownMigrationFiles(migrationDir string) error {
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

func handleErrorWithRemoveMigration(migrationDir string, err string) {
	color.Red(err)
	removeMigrationSetup(migrationDir)
	os.Exit(1)
}
