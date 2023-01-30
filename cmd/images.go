package cmd

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/config"
	"github.com/FasterBetter/cli-o-mat/util"
)

func showImages(shortHashes bool, images []*ec2.Image) {
	hashLen := 40
	if shortHashes {
		hashLen = 7
	}

	tableData := make([][]string, len(images))

	for idx, image := range images {
		var commit string

		for _, tag := range image.Tags {
			if aws.StringValue(tag.Key) == config.CommitTag {
				commit = aws.StringValue(tag.Value)

				if shortHashes && commit != "" {
					commit = commit[0:hashLen]
				}

				break
			}
		}

		tableData[idx] = []string{
			aws.StringValue(image.ImageId),
			aws.StringValue(image.Architecture),
			aws.StringValue(image.State),
			commit,
			aws.StringValue(image.Name),
			aws.StringValue(image.CreationDate),
		}
	}

	sort.Slice(tableData, func(i, j int) bool {
		return tableData[i][5] < tableData[j][5]
	})

	tableConfig := &util.Table{
		Columns: []util.Column{
			{Name: "ID"},
			{Name: "Arch"},
			{Name: "State"},
			{Name: "Commit"},
			{Name: "Name"},
			{Name: "Created"},
		},
	}

	tableConfig.Show(tableData)
}

// nolint: gochecknoglobals
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "List AMIs.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		omat := loadOmatConfig()

		details := awsutil.FindAndAssumeAdminRole(omat.BuildAccountSlug, omat)

		ec2Client := ec2.New(details.Session, details.Config)

		imageList, err := ec2Client.DescribeImages(&ec2.DescribeImagesInput{
			Owners: []*string{aws.String("self")},
		})
		if err != nil {
			util.Fatal(err)
		}

		showImages(imagesShortHashes, imageList.Images)
	},
}

// nolint: gochecknoglobals
var imagesShortHashes bool

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(imagesCmd)

	imagesCmd.Flags().BoolVarP(&imagesShortHashes, "short", "", false, "Shorten git commit SHAs")
}
