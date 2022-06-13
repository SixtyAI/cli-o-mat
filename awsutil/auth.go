package awsutil

import (
	"fmt"
	"strings"

	"github.com/FasterBetter/cli-o-mat/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/cockroachdb/errors"
)

func LoadCredentialsAndAssumeRole(omatConfig *config.Omat) (*session.Session, *aws.Config, error) {
	paramPrefix := omatConfig.Prefix()

	rootSess := session.Must(session.NewSession())
	ssmClient := ssm.New(rootSess, aws.NewConfig().WithRegion(omatConfig.Region))
	roleParamName := fmt.Sprintf("%s/ci-cd/roles/admin", paramPrefix)

	roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(roleParamName),
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "ParameterNotFound") {
			return nil, nil, errors.Errorf("couldn't find role parameter: %s", roleParamName)
		}

		return nil, nil, errors.WithStack(err)
	}

	arn := aws.StringValue(roleParam.Parameter.Value)
	if arn == "" {
		return nil, nil, errors.Errorf("paramater '%s' was empty", roleParamName)
	}

	assumedSess := session.Must(session.NewSession())
	creds := stscreds.NewCredentials(assumedSess, arn)

	return assumedSess, &aws.Config{Credentials: creds}, nil
}
