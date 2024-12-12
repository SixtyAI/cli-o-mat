package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/SixtyAI/cli-o-mat/config"
	"github.com/SixtyAI/cli-o-mat/util"
)

const (
	AWSAPIError = 13
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

func loadOmatConfig(accountName string) *config.Omat {
	omat := config.NewOmat(accountName)

	if err := omat.LoadConfig(); err != nil {
		util.Fatalf(1, "Failed to load omat config.\n")
	}

	if region != "" {
		omat.Region = region
	}

	if environment != "" {
		omat.Environment = environment
	}

	if deployService != "" {
		omat.DeployService = deployService
	}

	return omat
}
