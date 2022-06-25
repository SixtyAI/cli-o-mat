package cmd

import (
	"fmt"
	"time"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

func showHosts(hosts []*ec2.Instance) {
	// N.B. It's a bit rude to modify our parameter in-place in a function the caller expects to just
	// show things, but since we're not reusing the data structure yet, I'm not gonna add the overhead
	// of the slice copy.
	// sort.Slice(hosts, func(i, j int) bool {
	// 	return aws.StringValue(hosts[i].Name) < aws.StringValue(hosts[j].Name)
	// })

	// Default to be wide enough for the header.
	maxTypeLength := 4
	maxStateLength := 5
	maxApplicationLength := 12
	maxServiceLength := 8
	maxASGLength := 3

	for _, host := range hosts {
		if len(aws.StringValue(host.InstanceType)) > maxTypeLength {
			maxTypeLength = len(aws.StringValue(host.InstanceType))
		}

		state := host.State
		if state != nil {
			stateName := aws.StringValue(state.Name)
			if len(stateName) > maxStateLength {
				maxStateLength = len(stateName)
			}
		}

		for _, tag := range host.Tags {
			key := aws.StringValue(tag.Key)
			valueLen := len(aws.StringValue(tag.Value))
			if key == "Application" && valueLen > maxApplicationLength {
				maxApplicationLength = valueLen
			} else if key == "Service" && valueLen > maxServiceLength {
				maxServiceLength = valueLen
			} else if key == "aws:autoscaling:groupName" && valueLen > maxASGLength {
				maxASGLength = valueLen
			}
		}
	}

	hostsFormat := fmt.Sprintf("%%-19s %%-20s %%-%ds %%-7s %%-21s %%-%ds %%-15s %%-%ds %%-%ds %%4s %%-%ds %%s\n",
		maxTypeLength, maxStateLength, maxApplicationLength, maxServiceLength, maxASGLength)

	fmt.Printf(hostsFormat, "ID", "Launched", "Type", "Arch", "Image", "State", "Public IP",
		"Application", "Service", "Ver.", "ASG", "Keypair")

	for _, host := range hosts {
		var application string
		var service string
		var asg string
		var launchTemplateVersion string

		for _, tag := range host.Tags {
			key := aws.StringValue(tag.Key)
			val := aws.StringValue(tag.Value)
			if key == "Application" {
				application = val
			} else if key == "Service" {
				service = val
			} else if key == "aws:autoscaling:groupName" {
				asg = val
			} else if key == "aws:ec2launchtemplate:version" {
				launchTemplateVersion = val
			}
		}

		// host.IamInstanceProfile
		// host.SecurityGroups

		launchedAt := host.LaunchTime.Format(time.RFC3339)

		state := host.State
		var stateName string
		if state != nil {
			stateName = aws.StringValue(state.Name)
		}
		fmt.Printf(hostsFormat, aws.StringValue(host.InstanceId), launchedAt,
			aws.StringValue(host.InstanceType), aws.StringValue(host.Architecture),
			aws.StringValue(host.ImageId), stateName, aws.StringValue(host.PublicIpAddress), application,
			service, launchTemplateVersion, asg, aws.StringValue(host.KeyName))
	}
}

// nolint: gochecknoglobals
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List AMIs.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		omat, err := loadOmatConfig()
		if err != nil {
			util.Fatal(err)
		}

		details, err := awsutil.FindAndAssumeAdminRole(omat.DeployAccountSlug, omat)
		if err != nil {
			util.Fatal(err)
		}

		ec2Client := ec2.New(details.Session, details.Config)
		nextToken := aws.String("")
		hosts := []*ec2.Instance{}

		for {
			hostResp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
				NextToken: nextToken,
			})
			if err != nil {
				util.Fatal(err)
			}

			for _, res := range hostResp.Reservations {
				hosts = append(hosts, res.Instances...)
			}

			if hostResp.NextToken == nil {
				break
			}

			nextToken = hostResp.NextToken
		}

		showHosts(hosts)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(hostsCmd)
}
