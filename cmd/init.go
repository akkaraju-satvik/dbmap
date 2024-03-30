package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/akkaraju-satvik/dbmigo/queries"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	_ "github.com/lib/pq"
)

const configFileContent = `
{
	"migrations_dir": "$MIGRATIONS_DIR",
	"db_connection": "$DB_URL"
}
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dbmigo in a project directory",
	Long: `Initializes dbmigo in a project directory by creating a migrations directory and a config file.

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
			color.Red("Please provide a database connection string")
			os.Exit(1)
		}
		// check if connection is a postgres connection string
		postgresConnStringRegex := regexp.MustCompile(`^postgres:\/\/.*\:.*@.*\/.*$`)
		if !postgresConnStringRegex.MatchString(connection) {
			color.Red("Please provide a valid postgres connection string\n Example: postgres://user:password@localhost:5432/dbname\n")
			os.Exit(1)
		}
		ssl, _ := cmd.Flags().GetBool("ssl")
		if !strings.Contains(connection, "sslmode") && !ssl {
			connection = connection + "?sslmode=disable"
		}
		// create migrations directory
		err = os.MkdirAll(migrationsDir, 0755)
		if err != nil {
			color.Red("Error creating migrations directory\n %s", err.Error())
			os.Exit(1)
		}
		// create config file
		configFile, err := os.Create("dbmigo.json")
		if err != nil {
			color.Red("Error creating config file")
			os.Exit(1)
		}
		defer configFile.Close()
		configFileContent := strings.Replace(configFileContent, "$MIGRATIONS_DIR", migrationsDir, 1)
		configFileContent = strings.Replace(configFileContent, "$DB_URL", connection, 1)
		_, err = configFile.WriteString(configFileContent)
		if err != nil {
			color.Red("Error writing to config file")
			removeSetup()
			os.Exit(1)
		}
		db, err := sql.Open("postgres", connection)
		if err != nil {
			color.Red("Could not find a database with the provided connection string\n %s", err.Error())
			removeSetup()
			os.Exit(1)
		}
		defer db.Close()
		_, err = db.Exec(queries.InitQuery)
		if err != nil {
			color.Red("Error creating migrations table\n %s", err.Error())
			removeSetup()
			os.Exit(1)
		}
		color.Green("Setup complete")
	},
}

func init() {
	initCmd.Flags().StringP("migrations-dir", "d", "migrations", "Directory to store migration files")
	initCmd.Flags().StringP("db-connection", "c", "", "Database connection string")
	initCmd.MarkFlagRequired("db-connection")
	initCmd.Flags().BoolP("ssl", "s", false, "Use SSL for database connection")
	rootCmd.AddCommand(initCmd)
}

func removeSetup() {
	os.Remove("dbmigo.json")
	os.Remove("migrations")
}
