package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/config"
	"github.com/FasterBetter/cli-o-mat/util"
)

var errNoLaunchTemplate = fmt.Errorf("no launch template versions found")

func showLaunchTemplateVersions(shortHashes bool, templates []*ec2.LaunchTemplateVersion, groupMap map[string]string,
	imageMap map[string]string,
) {
	hashLen := 40
	if shortHashes {
		hashLen = 7
	}

	sort.Slice(templates, func(i, j int) bool {
		return aws.Int64Value(templates[i].VersionNumber) < aws.Int64Value(templates[j].VersionNumber)
	})

	tableData := make([][]string, len(templates))

	for idx, version := range templates {
		data := version.LaunchTemplateData
		if data == nil {
			data = &ec2.ResponseLaunchTemplateData{}
		}

		securityGroupIds := stringifySecurityGroups(version, groupMap)

		commit := imageMap[aws.StringValue(data.ImageId)]
		if shortHashes && commit != "" {
			commit = commit[0:hashLen]
		}

		tableData[idx] = []string{
			fmt.Sprintf("%d", aws.Int64Value(version.VersionNumber)),
			awsutil.DefaultToString(version.DefaultVersion),
			version.CreateTime.Format(time.RFC3339),
			aws.StringValue(data.InstanceType),
			aws.StringValue(data.ImageId),
			commit,
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
			{Name: "Commit"},
			{Name: "Keypair"},
			{Name: "Security Groups"},
		},
	}

	tableConfig.Show(tableData)
}

func stringifySecurityGroups(version *ec2.LaunchTemplateVersion, groupMap map[string]string) string {
	securityGroupIds := ""
	for _, group := range version.LaunchTemplateData.SecurityGroupIds {
		securityGroupIds += ", " + groupMap[aws.StringValue(group)]
	}

	if len(securityGroupIds) > 0 {
		securityGroupIds = securityGroupIds[2:]
	}

	return securityGroupIds
}

func buildSecurityGroupMapping(ec2Client *ec2.EC2, versions []*ec2.LaunchTemplateVersion) (map[string]string, error) {
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
		return nil, errors.WithStack(err)
	}

	for _, group := range groups {
		groupMap[aws.StringValue(group.GroupId)] = aws.StringValue(group.GroupName)
	}

	return groupMap, nil
}

func buildImageMapping(ec2Client *ec2.EC2) (map[string]string, error) {
	imageMap := map[string]string{}

	// N.B. I'd _like_ to just request the image IDs we care about, but the ImageId parameter here is... not behaving
	// nicely.  Specifically, if _any_ of the images are not found, the API will return an error.  Rather than try to
	// sort out what does/doesn't exist by parsing the error, we'll just request all of them.
	images, err := ec2Client.DescribeImages(&ec2.DescribeImagesInput{
		// ImageIds:          aws.StringSlice(imageIDs),
		IncludeDeprecated: aws.Bool(true),
		Owners:            []*string{aws.String("self")},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, image := range images.Images {
		var commit string

		for _, tag := range image.Tags {
			if aws.StringValue(tag.Key) == config.CommitTag {
				commit = aws.StringValue(tag.Value)

				break
			}
		}

		imageMap[aws.StringValue(image.ImageId)] = commit
	}

	return imageMap, nil
}

//nolint: gochecknoglobals
var launchTemplateCmd = &cobra.Command{
	Use:   "template template-name",
	Short: "Show details about a launch template.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		omat := loadOmatConfig()

		deployAcctDetails := awsutil.FindAndAssumeAdminRole(omat.DeployAccountSlug, omat)
		deployAcctEC2Client := ec2.New(deployAcctDetails.Session, deployAcctDetails.Config)

		buildAcctDetails := awsutil.FindAndAssumeAdminRole(omat.BuildAccountSlug, omat)
		buildAcctEC2Client := ec2.New(buildAcctDetails.Session, buildAcctDetails.Config)

		versions, err := awsutil.FetchLaunchTemplateVersions(deployAcctEC2Client, args[0], nil)
		if err != nil {
			util.Fatal(err)
		}

		if len(versions) == 0 {
			util.Fatal(errNoLaunchTemplate)
		}

		groupMap, err := buildSecurityGroupMapping(deployAcctEC2Client, versions)
		if err != nil {
			util.Fatal(err)
		}

		imageMap, err := buildImageMapping(buildAcctEC2Client)
		if err != nil {
			util.Fatal(err)
		}

		showLaunchTemplateVersions(templateShortHashes, versions, groupMap, imageMap)
	},
}

//nolint: gochecknoglobals
var templateShortHashes bool

//nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(launchTemplateCmd)

	launchTemplateCmd.Flags().BoolVarP(&templateShortHashes, "short", "", false, "Shorten git commit SHAs")
}
