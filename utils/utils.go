package utils

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"

	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

const ConfigFile = `migrations_dir: $MIGRATIONS_DIR
db_url: $DB_URL
`

type Config struct {
	MigrationsDir string `json:"migrations_dir" yaml:"migrations_dir"`
	DbURL         string `json:"db_url" yaml:"db_url"`
}

func ReadConfig(configFilePath string) (Config, error) {
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}
	if fileContent == nil {
		return Config{}, nil
	}
	var config Config
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

func HandleEmptyValuesInConfig(config Config) {
	var missingValues []string
	if config.MigrationsDir == "" {
		missingValues = append(missingValues, "migrations_dir")
	}
	if config.DbURL == "" {
		missingValues = append(missingValues, "db_url")
	}
	if len(missingValues) > 0 {
		color.Red("Please provide the following values in the config file: %s", strings.Join(missingValues, ", "))
		os.Exit(1)
	}
}

func GetUpMigration(migrationDir string) string {
	upMigration, err := os.ReadFile(migrationDir + "/up.sql")
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("up.sql file not found in migration directory")
			os.Exit(1)
		}
		color.Red("Error reading up migration file\n %s", err.Error())
		os.Exit(1)
	}
	return string(upMigration)
}

func GetDownMigration(migrationDir string) string {
	downMigration, err := os.ReadFile(migrationDir + "/down.sql")
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("down.sql file not found in migration directory")
			os.Exit(1)
		}
		color.Red("Error reading down migration file\n %s", err.Error())
		os.Exit(1)
	}
	return string(downMigration)
}

func ApplyMigration(dbUrl string, migration map[string]string) error {

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}

	_, err = conn.Exec(migration["migration"])
	if err != nil {
		return err
	}

	return nil
}

func GetMigrationQuery(migrationsDir, migrationName, migrationType string) (string, error) {
	migrationDir := migrationsDir + "/" + migrationName
	var migration string
	if migrationType == "up" {
		migration = GetUpMigration(migrationDir)
	} else {
		migration = GetDownMigration(migrationDir)
	}
	return migration, nil
}

func UpdateMigrationStatus(migrationId, migrationQuery string) error {
	conn, err := sql.Open("postgres", migrationQuery)
	if err != nil {
		return err
	}
	_, err = conn.Exec(queries.UpdateMigrationStatus, migrationId)
	if err != nil {
		return err
	}
	return nil
}
