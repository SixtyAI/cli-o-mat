package cmd

import (
	"github.com/spf13/cobra"
)

//nolint: gochecknoglobals
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Umbrella command for all things deploy_o_mat related.",
	Long:  ``,
}

//nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(deployCmd)
}
