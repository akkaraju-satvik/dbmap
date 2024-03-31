package utils

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

const ConfigFile = `migrations_dir: $MIGRATIONS_DIR
db_connection: $DB_URL
`

type Config struct {
	MigrationsDir string `json:"migrations_dir" yaml:"migrations_dir"`
	DBConnection  string `json:"db_connection" yaml:"db_connection"`
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
	if config.DBConnection == "" {
		missingValues = append(missingValues, "db_connection")
	}
	if len(missingValues) > 0 {
		color.Red("Please provide the following values in the config file: %s", strings.Join(missingValues, ", "))
		os.Exit(1)
	}
}
