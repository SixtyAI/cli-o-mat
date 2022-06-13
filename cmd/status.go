package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the deployment status of all services.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("status called")
	},
}

// nolint: gochecknoinits
func init() {
	deployCmd.AddCommand(statusCmd)
}
