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

func getTagValues(tags []*ec2.Tag) (string, string, string, string) {
	var (
		application           string
		service               string
		asg                   string
		launchTemplateVersion string
	)

	for _, tag := range tags {
		key := aws.StringValue(tag.Key)
		val := aws.StringValue(tag.Value)

		switch key {
		case "Application":
			application = val
		case "Service":
			service = val
		case "aws:autoscaling:groupName":
			asg = val
		case "aws:ec2launchtemplate:version":
			launchTemplateVersion = val
		}
	}

	return application, service, asg, launchTemplateVersion
}

func showHosts(hosts []*ec2.Instance) {
	tableData := make([][]string, len(hosts))

	for idx, host := range hosts {
		application, service, asg, launchTemplateVersion := getTagValues(host.Tags)

		var stateName string

		state := host.State
		if state != nil {
			stateName = aws.StringValue(state.Name)
		}

		tableData[idx] = []string{
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
