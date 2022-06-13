package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var manualCmd = &cobra.Command{
	Use:   "manual",
	Short: "Shows the full documentation for this tool",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`Omat Manual

Configuration

omat will look for a file named .omat.yml, starting in the current directory and
proceeding up the file system.  This file is required, although all options can
also be specified with environment variables, and some via CLI flags (run this
tool with -h for details)

+--------------------+--------------------------+------------------+
| YAML Key           | Environment Variable     | Flag             |
+--------------------+--------------------------+------------------+
| organizationPrefix | OMAT_ORGANIZATION_PREFIX |                  |
| region             | OMAT_REGION              | --region         |
| environment        | OMAT_ENVIRONMENT         | --env            |
| deployService      | OMAT_DEPLOY_SERVICE      | --deploy-service |
| buildAccountSlug   | OMAT_BUILD_ACCOUNT_SLUG  | --build-slug     |
+--------------------+--------------------------+------------------+

Run without options for a list of sub-commands.`)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(manualCmd)
}
