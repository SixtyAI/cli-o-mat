/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: stylecheck
//go:embed .version_string
var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cli-o-mat v%s\n", Version) // nolint: forbidigo
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
