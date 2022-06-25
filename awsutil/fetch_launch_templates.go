package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func FetchLaunchTemplates(ec2Client *ec2.EC2) ([]*ec2.LaunchTemplate, error) {
	nextToken := aws.String("")
	templates := []*ec2.LaunchTemplate{}

	for {
		resp, err := ec2Client.DescribeLaunchTemplates(&ec2.DescribeLaunchTemplatesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to describe launch templates")
		}

		templates = append(templates, resp.LaunchTemplates...)

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
	}

	return templates, nil
}
