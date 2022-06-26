package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

var errNoLaunchTemplate = fmt.Errorf("no launch template versions found")

func showLaunchTemplateVersions(templates []*ec2.LaunchTemplateVersion, groupMap map[string]string) {
	sort.Slice(templates, func(i, j int) bool {
		return aws.Int64Value(templates[i].VersionNumber) < aws.Int64Value(templates[j].VersionNumber)
	})

	tableData := make([][]string, len(templates))

	for idx, version := range templates {
		var isDefault string
		if aws.BoolValue(version.DefaultVersion) {
			isDefault = "yes"
		} else {
			isDefault = "" // Blank so the `yes` value stands out more.
		}

		data := version.LaunchTemplateData
		if data == nil {
			data = &ec2.ResponseLaunchTemplateData{}
		}

		securityGroupIds := ""
		for _, group := range version.LaunchTemplateData.SecurityGroupIds {
			securityGroupIds += ", " + groupMap[aws.StringValue(group)]
		}

		if len(securityGroupIds) > 0 {
			securityGroupIds = securityGroupIds[2:]
		}

		tableData[idx] = []string{
			fmt.Sprintf("%d", aws.Int64Value(version.VersionNumber)),
			isDefault,
			version.CreateTime.Format(time.RFC3339),
			aws.StringValue(data.InstanceType),
			aws.StringValue(data.ImageId),
			aws.StringValue(data.KeyName),
			securityGroupIds,
		}
	}

	tableConfig := &util.Table{
		Columns: []util.Column{
			{Name: "Version", RightAlign: true},
			{Name: "Default?"},
			{Name: "Created"},
			{Name: "Type"},
			{Name: "Image"},
			{Name: "Keypair"},
			{Name: "Security Groups"},
		},
	}

	tableConfig.Show(tableData)
}

// nolint: gochecknoglobals
var launchTemplateCmd = &cobra.Command{
	Use:   "template template-name",
	Short: "Show details about a launch template.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
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

		versions, err := awsutil.FetchLaunchTemplateVersions(ec2Client, args[0], nil)
		if err != nil {
			util.Fatal(err)
		}

		if len(versions) == 0 {
			util.Fatal(errNoLaunchTemplate)
		}

		groupMap := map[string]string{}
		for _, version := range versions {
			for _, groupID := range version.LaunchTemplateData.SecurityGroupIds {
				groupMap[aws.StringValue(groupID)] = ""
			}
		}

		groupIDs := make([]string, 0, len(groupMap))
		for groupID := range groupMap {
			groupIDs = append(groupIDs, groupID)
		}

		groups, err := awsutil.FetchSecurityGroups(ec2Client, groupIDs)
		if err != nil {
			util.Fatal(err)
		}

		for _, group := range groups {
			groupMap[aws.StringValue(group.GroupId)] = aws.StringValue(group.GroupName)
		}

		fmt.Printf("%+v\n", groups)
		showLaunchTemplateVersions(versions, groupMap)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(launchTemplateCmd)
}
