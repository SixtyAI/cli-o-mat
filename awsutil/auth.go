package awsutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/FasterBetter/cli-o-mat/config"
	"github.com/FasterBetter/cli-o-mat/util"
)

const (
	CantFindRoleParam       = 10
	ErrorLookingUpRoleParam = 11
	RoleParamEmpty          = 12
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
			util.Fatalf(CantFindRoleParam, "Couldn't find role parameter: %s\n", roleParamName)
		}

		fmt.Printf("Error looking up role parameter %s, got: %s\n", roleParamName, err.Error())
		util.Fatalf(ErrorLookingUpRoleParam, "Error looking up role parameter %s, got: %s\n", roleParamName, err.Error())
	}

	arn := aws.StringValue(roleParam.Parameter.Value)
	if arn == "" {
		util.Fatalf(RoleParamEmpty, "Paramater '%s' was empty.\n", roleParamName)
	}

	details := omat.Credentials.ForARN(arn)

	return details
}
