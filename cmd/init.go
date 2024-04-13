package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/akkaraju-satvik/dbmap/config"
	"github.com/akkaraju-satvik/dbmap/queries"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

var initErrorStrings = map[string]string{
	"migrations": "Error creating migrations directory\n %s",
	"config":     "Error creating config file",
	"table":      "Error creating migrations table\n %s",
	"db":         "Could not find a database with the provided connection string\n %s",
	"postgres":   "Please provide a valid postgres connection string\n Example: postgres://user:password@localhost:5432/dbname\n",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dbmap in a project directory",
	Long: `Initializes dbmap in a project directory by creating a migrations directory and a config file.

The config file will contain the database connection string and the migrations directory path.`,
	Run: func(cmd *cobra.Command, args []string) {
		migrationsDir, err := cmd.Flags().GetString("migrations-dir")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		connection, err := cmd.Flags().GetString("db-connection")
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		if connection == "" {
			fmt.Printf("Enter the database connection string: ")
			fmt.Scanln(&connection)
			if connection == "" {
				color.Red("Please provide a valid postgres connection string\n Example: postgres://user:password@localhost:5432/dbname\n")
				os.Exit(1)
			}
		}
		ssl, _ := cmd.Flags().GetBool("ssl")
		if !strings.Contains(connection, "sslmode=require") && !ssl {
			connection = connection + "?sslmode=disable"
		}
		err = initializeProject(migrationsDir, connection)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
	},
}

func initializeProject(migrationsDir, connection string) error {
	postgresConnStringRegex := regexp.MustCompile(`^postgres:\/\/.*\:.*@.*\/.*$`)
	if !postgresConnStringRegex.MatchString(connection) {
		return fmt.Errorf(initErrorStrings["postgres"])
	}
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	err := os.MkdirAll(migrationsDir, 0755)
	if err != nil {
		removeSetup()
		return fmt.Errorf(initErrorStrings["migrations"], err.Error())
	}
	configFile, err := os.Create("dbmap.config.yaml")
	if err != nil {
		removeSetup()
		return fmt.Errorf(initErrorStrings["config"])
	}
	defer configFile.Close()
	configFileContent := strings.Replace(config.ConfigFile, "$MIGRATIONS_DIR", migrationsDir, 1)
	configFileContent = strings.Replace(configFileContent, "$DB_URL", connection, 1)
	_, err = configFile.WriteString(configFileContent)
	if err != nil {
		removeSetup()
		return fmt.Errorf(initErrorStrings["config"])
	}
	db, err := sql.Open("postgres", connection)
	if err != nil {
		removeSetup()
		return fmt.Errorf(initErrorStrings["db"], err.Error())
	}
	defer db.Close()
	_, err = db.Exec(queries.InitQuery)
	if err != nil {
		removeSetup()
		return fmt.Errorf(initErrorStrings["table"], err.Error())
	}
	color.Green("Setup complete")
	return nil
}

func init() {
	initCmd.Flags().StringP("migrations-dir", "d", "migrations", "Directory to store migration files")
	initCmd.Flags().StringP("db-connection", "c", "", "Database connection string")
	initCmd.Flags().BoolP("ssl", "s", false, "Use SSL for database connection")
	rootCmd.AddCommand(initCmd)
}

func removeSetup() {
	os.Remove("dbmap.config.yaml")
	os.Remove("migrations")
}
