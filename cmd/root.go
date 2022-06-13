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
	rootCmd.PersistentFlags().StringVarP(&region, "region", "", "", "Which AWS region to operate in")
	rootCmd.PersistentFlags().StringVarP(&environment, "env", "", "", "Which logical environment to operate in")
	rootCmd.PersistentFlags().StringVarP(&deployService, "deploy-service", "", "", "The name of the deploy_o_mat service")
}

// nolint: gochecknoglobals
var (
	region        string
	environment   string
	deployService string
)

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

	if region != "" {
		omatConfig.Region = region
	}

	if environment != "" {
		omatConfig.Environment = environment
	}

	if deployService != "" {
		omatConfig.DeployService = deployService
	}

	return omatConfig, nil
}
