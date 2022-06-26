package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func FetchLaunchTemplateVersions(ec2Client *ec2.EC2, name string) ([]*ec2.LaunchTemplateVersion, error) {
	nextToken := aws.String("")
	versions := []*ec2.LaunchTemplateVersion{}

	for {
		resp, err := ec2Client.DescribeLaunchTemplateVersions(&ec2.DescribeLaunchTemplateVersionsInput{
			LaunchTemplateName: aws.String(name),
			NextToken:          nextToken,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to describe launch template versions")
		}

		versions = append(versions, resp.LaunchTemplateVersions...)

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
	}

	return versions, nil
}
