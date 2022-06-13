package awsutil

import (
	"fmt"
	"strings"

	"github.com/FasterBetter/cli-o-mat/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cockroachdb/errors"
)

func FindAndAssumeAdminRole(accountSlug string, omat *config.Omat) (*config.SessionDetails, error) {
	paramPrefix := omat.Prefix()

	ssmClient := ssm.New(omat.Credentials.RootSession, omat.Credentials.RootAWSConfig)
	roleParamName := fmt.Sprintf("%s/%s/roles/admin", paramPrefix, accountSlug)
	fmt.Printf("Looking for SSM parameter %s\n", roleParamName)

	roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(roleParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			return nil, errors.Errorf("couldn't find role parameter: %s", roleParamName)
		}

		return nil, errors.WithStack(err)
	}

	arn := aws.StringValue(roleParam.Parameter.Value)
	if arn == "" {
		return nil, errors.Errorf("paramater '%s' was empty", roleParamName)
	}

	details := omat.Credentials.ForARN(arn)

	return details, nil
}
