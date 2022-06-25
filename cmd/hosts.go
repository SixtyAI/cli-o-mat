package cmd

import (
	"sort"
	"time"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

func showHosts(hosts []*ec2.Instance) {
	tableData := make([][]string, len(hosts))

	for i, host := range hosts {
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

		state := host.State
		var stateName string
		if state != nil {
			stateName = aws.StringValue(state.Name)
		}

		tableData[i] = []string{
			aws.StringValue(host.InstanceId), host.LaunchTime.Format(time.RFC3339),
			aws.StringValue(host.InstanceType), aws.StringValue(host.Architecture),
			aws.StringValue(host.ImageId), stateName, aws.StringValue(host.PublicIpAddress), application,
			service, launchTemplateVersion, asg, aws.StringValue(host.KeyName),
		}
	}

	sort.Slice(tableData, func(i, j int) bool {
		return tableData[i][1] < tableData[j][1]
	})

	tableConfig := &util.Table{
		Columns: []util.Column{
			{Name: "ID"},
			{Name: "Launched"},
			{Name: "Type"},
			{Name: "Arch"},
			{Name: "Image"},
			{Name: "State"},
			{Name: "Public IP"},
			{Name: "Application"},
			{Name: "Service"},
			{Name: "Ver.", RightAlign: true},
			{Name: "ASG"},
			{Name: "Keypair"},
		},
	}

	tableConfig.Show(tableData)
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
