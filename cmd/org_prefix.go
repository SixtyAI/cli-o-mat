package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	CantFindOrgPrefix       = 14
	ErrorLookingUpRoleParam = 15
	OrgPrefixEmpty          = 16
)

// nolint: gochecknoglobals
var orgPrefixCmd = &cobra.Command{
	Use:   "org-prefix",
	Short: "Show the detected org prefix.",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		omat := loadOmatConfig("") // TODO: Fixme?

		fmt.Printf("Organization prefix: %s\n", omat.OrganizationPrefix)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(orgPrefixCmd)
}
