package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const ConfigFile = `migrations_dir: $MIGRATIONS_DIR
db_url: $DB_URL
`

type Config struct {
	MigrationsDir string `json:"migrations_dir" yaml:"migrations_dir"`
	DbURL         string `json:"db_url" yaml:"db_url"`
}

type MissingValuesError struct {
	MissingValues []string
}

func Read(configFilePath string) (Config, error) {
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
	missingValues := handleEmptyValues(config)
	if missingValues != nil {
		return config, errors.New("missing values in config file: " + strings.Join(missingValues, ", "))
	}
	return config, nil
}

func handleEmptyValues(config Config) []string {
	var missingValues []string
	if config.MigrationsDir == "" {
		missingValues = append(missingValues, "migrations_dir")
	}
	if config.DbURL == "" {
		missingValues = append(missingValues, "db_url")
	}
	if len(missingValues) > 0 {
		return missingValues
	}
	return nil
}
