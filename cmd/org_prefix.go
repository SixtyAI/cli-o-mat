package cmd

import (
	"fmt"
	"strings"

	"github.com/SixtyAI/cli-o-mat/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	Run: func(cmd *cobra.Command, args []string) {
		omat := loadOmatConfig()

		ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
		roleParamName := "/omat/organization_prefix"
		fmt.Printf("Looking for SSM parameter %s\n", roleParamName)

		roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
			Name: aws.String(roleParamName),
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "ParameterNotFound") {
				util.Fatalf(CantFindOrgPrefix, "Couldn't find org prefix parameter: %s\n", roleParamName)
			}

			fmt.Printf("Error looking up org prefix parameter %s, got: %s\n", roleParamName, err.Error())
			util.Fatalf(ErrorLookingUpRoleParam, "Error looking up org prefix parameter %s, got: %s\n", roleParamName, err.Error())
		}

		orgPrefix := aws.StringValue(roleParam.Parameter.Value)
		if orgPrefix == "" {
			util.Fatalf(OrgPrefixEmpty, "Paramater '%s' was empty.\n", roleParamName)
		}

		fmt.Printf("Organization prefix: %s\n", orgPrefix)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(orgPrefixCmd)
}
