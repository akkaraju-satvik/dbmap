package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
)

func TestInitializeProject(t *testing.T) {
	// Set up test data

	dir, err := os.MkdirTemp("", "test_init")
	if err != nil {
		t.Fatalf("Error creating temp directory: %s", err)
	}

	defer os.RemoveAll(dir)

	if err = os.Chdir(dir); err != nil {
		t.Fatalf("Error changing directory: %s", err)
	}

	pool, resource, port := createResource()
	defer pool.Purge(resource)

	migrationsDir := "test_migrations"
	connection := "postgres://postgres:postgres@localhost:" + port + "/postgres?sslmode=disable"

	// Call the function being tested
	err = initializeProject(migrationsDir, connection)
	if err != nil {
		t.Fatalf("Error initializing project: %s", err)
	}

}

func TestInitializeProjectError(t *testing.T) {
	tests := []struct {
		name          string
		migrationsDir string
		test_id       string
		connection    string
	}{
		{
			name:          "Invalid connection string",
			test_id:       "invalid_connection_string",
			migrationsDir: "test_migrations",
			connection:    "invalid_connection_string",
		},
		{
			name:          "Empty connection string",
			test_id:       "empty_connection_string",
			migrationsDir: "test_migrations",
			connection:    "",
		},
		{
			name:          "SSL error",
			test_id:       "ssl_error",
			migrationsDir: "test_migrations",
			connection:    "postgres://postgres:postgres@localhost:5432/postgres?sslmode=require",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("Error creating temp directory: %s", err)
			}

			defer os.RemoveAll(dir)

			if err = os.Chdir(dir); err != nil {
				t.Fatalf("Error changing directory: %s", err)
			}

			migrationsDir := tt.migrationsDir
			connection := tt.connection

			err = initializeProject(migrationsDir, connection)
			if err == nil {
				t.Fatalf("Expected error initializing project: %s", err)
			}
		})
	}
}

func createResource() (*dockertest.Pool, *dockertest.Resource, string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}
	resource, err := pool.Run("postgres", "15", []string{"POSTGRES_PASSWORD=postgres", "POSTGRES_USER=postgres", "POSTGRES_DB=postgres"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	port := resource.GetPort("5432/tcp")
	pool.MaxWait = 30 * time.Second
	if err = pool.Retry(func() error {
		dockerDB, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres?sslmode=disable", port))
		if err != nil {
			return err
		}
		defer dockerDB.Close()
		return dockerDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource, port
}
