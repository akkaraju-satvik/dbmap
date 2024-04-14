package cmd

import (
	"os"
	"testing"

	"github.com/akkaraju-satvik/dbmap/config"
)

func TestCreateMigration(t *testing.T) {
	tests := []struct {
		name        string
		testID      string
		expectedErr bool
		config      config.Config
	}{
		{
			name:        "Create Migration",
			testID:      "create_migration_success",
			expectedErr: false,
			config: config.Config{
				MigrationsDir: "migrations",
				DbURL:         "",
			},
		},
		{
			name:        "Invalid Config with Empty DB URL Value",
			testID:      "invalid_config_empty_db_url",
			expectedErr: true,
			config: config.Config{
				MigrationsDir: "migrations",
				DbURL:         "",
			},
		},
		{
			name:        "Invalid Config with Empty Migrations Directory Value",
			testID:      "invalid_config_empty_migrations_dir",
			expectedErr: true,
			config: config.Config{
				MigrationsDir: "",
				DbURL:         "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "test*")
			if err != nil {
				t.Fatalf("Error creating temp dir\n %s", err)
			}
			defer os.RemoveAll(dir)
			if err = os.Chdir(dir); err != nil {
				t.Fatalf("Error changing directory: %s", err)
			}
			pool, resource, port := createResource()
			defer pool.Purge(resource)
			dbUrl := "postgres://postgres:postgres@localhost:" + port + "/postgres?sslmode=disable"
			err = initializeProject("migrations", dbUrl)
			if err != nil {
				t.Fatalf("Error initializing a test project\n %s", err)
			}
			if tt.testID != "invalid_config_empty_db_url" {
				tt.config.DbURL = dbUrl
			}
			err = createMigration(tt.config)
			if err != nil {
				if !tt.expectedErr {
					t.Fatalf("Error creating migration\n %s", err.Error())
				}
			}
		})
	}
}
