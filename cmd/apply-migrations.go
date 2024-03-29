package cmd

import (
	"github.com/akkaraju-satvik/dbmigo/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var applyMigrationsCmd = &cobra.Command{
	Use:   "apply-migrations",
	Short: "Apply all migrations in the migrations directory specified in the config file",
	Long: `Applies all migrations in the migrations directory specified in the config file.

The migrations will be applied in the order of their creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := cmd.Flag("config-file").Value.String()
		config, err := utils.ReadConfig(configFilePath)
		if err != nil {
			color.Red("Error reading config file\n %s", err.Error())
		}
		utils.HandleEmptyValuesInConfig(config)

	},
}

func init() {
	applyMigrationsCmd.Flags().StringP("config-file", "c", "dbmigo.json", "Path to the config file")
	rootCmd.AddCommand(applyMigrationsCmd)
}
