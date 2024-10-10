package cmd

import (
	_ "embed"
	"log"

	"github.com/spf13/cobra"
)

// nolint: stylecheck
//
//go:embed .version_string
var Version string

// nolint: gochecknoglobals
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version.",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		log.Printf("cli-o-mat v%s", Version)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(versionCmd)
}
