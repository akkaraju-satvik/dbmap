package migrations

import (
	"database/sql"
	"os"

	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
)

func GetUp(migrationDir string) string {
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

func GetDown(migrationDir string) string {
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

func Apply(dbUrl string, migration map[string]string) error {

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

func GetQuery(migrationsDir, migrationName, migrationType string) (string, error) {
	migrationDir := migrationsDir + "/" + migrationName
	var migration string
	if migrationType == "up" {
		migration = GetUp(migrationDir)
	} else {
		migration = GetDown(migrationDir)
	}
	return migration, nil
}

func UpdateStatus(migrationId, migrationQuery string) error {
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
