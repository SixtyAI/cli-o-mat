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

-------------
Configuration
-------------

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
| deployAccountSlug  | OMAT_DEPLOY_ACCOUNT_SLUG | --deploy-slug    |
+--------------------+--------------------------+------------------+

Run without options for a list of sub-commands.

-----------
Error Codes
-----------

The following are error codes that may be returned from multiple sub-commands.
See the help for each sub-command for details on what other errors it may
return.

1  - Invalid configuration.  Either the .omat.yml file is missing, or its
     contents are invalid.
10 - Couldn't find the SSM parameter specifying the name of the role to assume
     for admin access.
11 - Some other error occurred when fetching the admin role name from SSM.
12 - The SSM parameter specifying the name of the role to assume for admin
     access was found, but its value was empty.
13 - AWS API error.  This is a generic error code for any AWS API not
     specifically handled by any command.`)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(manualCmd)
}
