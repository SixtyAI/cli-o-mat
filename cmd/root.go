package cmd

import (
	"os"

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
