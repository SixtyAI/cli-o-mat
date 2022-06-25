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

func getLaunchTemplateTagValues(tags []*ec2.Tag) (string, string) {
	var (
		application string
		service     string
	)

	for _, tag := range tags {
		key := aws.StringValue(tag.Key)
		val := aws.StringValue(tag.Value)

		switch key {
		case "Application":
			application = val
		case "Service":
			service = val
		}
	}

	return application, service
}

func showLaunchTemplates(templates []*ec2.LaunchTemplate) {
	tableData := make([][]string, len(templates))

	for idx, template := range templates {
		application, service := getLaunchTemplateTagValues(template.Tags)

		tableData[idx] = []string{
			aws.StringValue(template.LaunchTemplateId),
			aws.StringValue(template.LaunchTemplateName),
			application,
			service,
			fmt.Sprintf("%d", aws.Int64Value(template.DefaultVersionNumber)),
			fmt.Sprintf("%d", aws.Int64Value(template.LatestVersionNumber)),
			template.CreateTime.Format(time.RFC3339),
		}
	}

	sort.Slice(tableData, func(i, j int) bool {
		return tableData[i][1] < tableData[j][1]
	})

	tableConfig := &util.Table{
		Columns: []util.Column{
			{Name: "ID"},
			{Name: "Name"},
			{Name: "Application"},
			{Name: "Service"},
			{Name: "Default", RightAlign: true},
			{Name: "Latest", RightAlign: true},
			{Name: "Created"},
		},
	}

	tableConfig.Show(tableData)
}

// nolint: gochecknoglobals
var launchTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List launch templates.",
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
		templates, err := awsutil.FetchLaunchTemplates(ec2Client, nil)
		if err != nil {
			util.Fatal(err)
		}

		showLaunchTemplates(templates)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(launchTemplatesCmd)
}
