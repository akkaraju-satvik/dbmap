package migrations

import (
	"os"

	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
	"github.com/jmoiron/sqlx"
)

type MigrationType int
type MigrationStatus int

const (
	UP MigrationType = iota
	DOWN
)

const (
	PENDING MigrationStatus = iota
	APPLIED
)

func (m MigrationType) String() string {
	return []string{"UP", "DOWN"}[m]
}

func (s MigrationStatus) String() string {
	return []string{"PENDING", "APPLIED"}[s]
}

type Migration struct {
	MigrationId     string `db:"migration_id"`
	MigrationName   string `db:"migration_name"`
	MigrationTime   string `db:"migration_time"`
	MigrationStatus string `db:"migration_status"`
}

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

func GetMigrations(dbUrl, migrationType string, migrationList *[]Migration) error {
	conn, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer conn.Close()
	var migrationStatus string
	if migrationType == UP.String() {
		migrationStatus = PENDING.String()
	} else {
		migrationStatus = APPLIED.String()
	}
	err = conn.Select(migrationList, queries.GetMigrations, migrationStatus)
	if err != nil {
		return err
	}
	return nil
}

func GetQuery(migrationsDir, migrationName, migrationType string) (string, error) {
	migrationDir := migrationsDir + "/" + migrationName
	var migration string
	if migrationType == UP.String() {
		migration = GetUp(migrationDir)
	} else if migrationType == DOWN.String() {
		migration = GetDown(migrationDir)
	}
	return migration, nil
}

func ApplyMigration(dbUrl, migrationQuery, migrationType string, migration Migration) error {
	conn, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer conn.Close()
	// Apply the migration
	_, err = conn.Exec(migrationQuery)
	if err != nil {
		return err
	}
	// Update Migration Status
	var migrationStatus string
	if migrationType == "DOWN" {
		migrationStatus = PENDING.String()
	} else {
		migrationStatus = APPLIED.String()
	}
	_, err = conn.Exec(queries.UpdateMigrationStatus, migrationStatus, migration.MigrationId)
	if err != nil {
		return err
	}
	// Insert Migration Query
	_, err = conn.Exec(queries.InsertMigrationQuery, migration.MigrationId, migrationQuery, migrationType)
	if err != nil {
		return err
	}
	return nil
}
