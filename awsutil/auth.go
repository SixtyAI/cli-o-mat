package awsutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/FasterBetter/cli-o-mat/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func FindAndAssumeAdminRole(accountSlug string, omat *config.Omat) *config.SessionDetails {
	paramPrefix := omat.Prefix()

	ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
	roleParamName := fmt.Sprintf("%s/%s/roles/admin", paramPrefix, accountSlug)
	fmt.Printf("Looking for SSM parameter %s\n", roleParamName)

	roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(roleParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			fmt.Printf("Couldn't find role parameter: %s\n", roleParamName)
			os.Exit(10)
		}

		fmt.Printf("Error looking up role parameter %s, got: %s\n", roleParamName, err)
		os.Exit(11)
	}

	arn := aws.StringValue(roleParam.Parameter.Value)
	if arn == "" {
		fmt.Printf("Paramater '%s' was empty\n", roleParamName)
		os.Exit(12)
	}

	details := omat.Credentials.ForARN(arn)

	return details
}
