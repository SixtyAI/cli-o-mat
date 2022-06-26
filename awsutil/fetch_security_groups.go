package awsutil

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cockroachdb/errors"
)

func FetchSecurityGroups(ec2Client *ec2.EC2, groupIDs []string) ([]*ec2.SecurityGroup, error) {
	var nextToken *string

	groups := []*ec2.SecurityGroup{}
	groupIDsRef := aws.StringSlice(groupIDs)

	for {
		resp, err := ec2Client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
			GroupIds:  groupIDsRef,
			NextToken: nextToken,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		groups = append(groups, resp.SecurityGroups...)

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
		groupIDsRef = nil
	}

	return groups, nil
}
