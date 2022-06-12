package cmd

import (
	"fmt"
	"os"

	"github.com/FasterBetter/cli-o-mat/config"

	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "cli-o-mat",
	Short: "CLI tool for managing Omat deploys",
	Long: `cli-o-mat is a tool for seeing what's deployable, what's deployed,
initiating deploys, and cancelling deploys using Teak.io's Omat infrastructure
tooling.`,
}

func Execute() {
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}

// nolint: gochecknoinits
func init() {
	rootCmd.PersistentFlags().StringVarP(&environment, "env", "", "", "Which environment to operate in")
}

var environment string // nolint: gochecknoglobals

func loadOmatConfig() (*config.Omat, error) {
	configFile, err := config.FindOmatConfig(".")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	omatConfig := config.NewOmat()

	if err = omatConfig.LoadConfigFromFile(configFile); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	omatConfig.LoadConfigFromEnv()

	if environment != "" {
		omatConfig.Environment = environment
	}

	return omatConfig, nil
}
