/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/akkaraju-satvik/dbmigo/queries"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type config struct {
	MigrationsDir string `json:"migrations_dir" yaml:"migrations_dir"`
	DBConnection  string `json:"db_connection" yaml:"db_connection"`
}

// createMigrationCmd represents the createMigration command
var createMigrationCmd = &cobra.Command{
	Use:   "create-migration",
	Short: "Create a new migration in the migrations directory specified in the config file",
	Long: `Creates a new migration in the migrations directory specified in the config file.

The migration will have an up.sql and down.sql file with placeholders for the migration queries.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config-file").Value.String()
		config, err := readConfig(configFilePath)
		if err != nil {
			color.Red("Error reading config file\n %s", err.Error())
			os.Exit(1)
		}
		if config.MigrationsDir == "" {
			color.Red("Please provide a migrations directory in the config file")
			os.Exit(1)
		}
		if config.DBConnection == "" {
			color.Red("Please provide a database connection string in the config file")
			os.Exit(1)
		}
		currentTime := fmt.Sprintf("%d", time.Now().Unix())
		migrationDir := config.MigrationsDir + "/" + currentTime
		if err = os.MkdirAll(migrationDir, 0755); err != nil {
			color.Red("Error creating migration directory\n %s", err.Error())
			os.Exit(1)
		}
		upFile, err := os.Create(migrationDir + "/up.sql")
		if err != nil {
			color.Red("Error creating migration file\n %s", err.Error())
			removeMigrationSetup(migrationDir)
			os.Exit(1)
		}
		fmt.Fprintf(upFile, "-- Write your UP migration here\n\n")
		upFile.Close()
		downFile, err := os.Create(migrationDir + "/down.sql")
		if err != nil {
			color.Red("Error creating migration file\n %s", err.Error())
			removeMigrationSetup(migrationDir)
			os.Exit(1)
		}
		fmt.Fprintf(downFile, "-- Write your DOWN migration here\n\n")
		downFile.Close()
		db, err := sql.Open("postgres", config.DBConnection)
		if err != nil {
			color.Red("Error connecting to the database\n %s", err.Error())
		}
		_, err = db.Exec(queries.CreateMigrationQuery, currentTime)
		if err != nil {
			color.Red("Error creating migration\n %s", err.Error())
			removeMigrationSetup(migrationDir)
			os.Exit(1)
		}
		db.Close()
		color.Green("Migration created successfully")
	},
}

func init() {
	createMigrationCmd.Flags().StringP("config-file", "c", "pgo-migrator.json", "config file to use")
	rootCmd.AddCommand(createMigrationCmd)
}

func readConfig(configFilePath string) (config, error) {
	// Read the config file
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return config{}, err
	}
	if fileContent == nil {
		return config{}, nil
	}
	var config config
	if strings.HasSuffix(configFilePath, ".yaml") || strings.HasSuffix(configFilePath, ".yml") {
		err = yaml.Unmarshal(fileContent, &config)
	} else {
		err = json.Unmarshal(fileContent, &config)
	}
	if err != nil {
		return config, err
	}
	return config, nil
}

func removeMigrationSetup(migrationDir string) {
	os.RemoveAll(migrationDir)
}
