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
	rootCmd.PersistentFlags().StringVarP(&buildAccountSlug, "build-slug", "", "", "The build-account slug (e.g. ci-cd)")
	rootCmd.PersistentFlags().StringVarP(&deployAccountSlug, "deploy-slug", "", "",
		"The deploy-account slug (e.g. workload)")
}

// nolint: gochecknoglobals
var (
	region            string
	environment       string
	deployService     string
	buildAccountSlug  string
	deployAccountSlug string
)

func loadOmatConfig() *config.Omat {
	omat := config.NewOmat()

	if err := omat.LoadConfig(); err != nil {
		fmt.Printf("Failed to load omat config")
		os.Exit(1)
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

	if buildAccountSlug != "" {
		omat.BuildAccountSlug = buildAccountSlug
	}

	if deployAccountSlug != "" {
		omat.DeployAccountSlug = deployAccountSlug
	}

	omat.InitCredentials()

	return omat
}
