package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func FetchLaunchTemplateVersions(ec2Client *ec2.EC2, name string,
	desiredVersion *string,
) ([]*ec2.LaunchTemplateVersion, error) {
	var nextToken *string

	var desiredVersions []*string
	if desiredVersion != nil {
		desiredVersions = []*string{desiredVersion}
	}

	versions := []*ec2.LaunchTemplateVersion{}

	for {
		resp, err := ec2Client.DescribeLaunchTemplateVersions(&ec2.DescribeLaunchTemplateVersionsInput{
			LaunchTemplateName: aws.String(name),
			NextToken:          nextToken,
			Versions:           desiredVersions,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		versions = append(versions, resp.LaunchTemplateVersions...)

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
		desiredVersions = nil
	}

	return versions, nil
}
