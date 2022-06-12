package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Shows the full documentation for this tool",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`omat
Configuration

omat will look for a file named .omat.yml, starting in the current directory and
proceeding up the file system.  This file is required, although all options can
also be specified with environment variables, and some via CLI flags (run this
tool with -h for details)

+--------------------+--------------------------+----------+
| YAML Key           | Environment Variable     | Flag     |
+--------------------+--------------------------+----------+
| organizationPrefix | OMAT_ORGANIZATION_PREFIX |          |
| region             | OMAT_REGION              | --region |
| environment        | OMAT_ENVIRONMENT         | --env    |
+--------------------+--------------------------+----------+

Run without options for a list of sub-commands.`)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(helpCmd)
}
