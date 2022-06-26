package awsutil

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cockroachdb/errors"
)

func FetchSubnets(ec2Client *ec2.EC2) ([]*ec2.Subnet, error) {
	var nextToken *string
	subnets := []*ec2.Subnet{}

	for {
		resp, err := ec2Client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		subnets = append(subnets, resp.Subnets...)

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
	}

	return subnets, nil
}
