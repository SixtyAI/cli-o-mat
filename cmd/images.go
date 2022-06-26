package cmd

import (
	"sort"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
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
			if aws.StringValue(tag.Key) == "BuildCommit" {
				commit = aws.StringValue(tag.Value)

				if shortHashes {
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
		}
	}

	sort.Slice(tableData, func(i, j int) bool {
		return tableData[i][4] < tableData[j][4]
	})

	tableConfig := &util.Table{
		Columns: []util.Column{
			{Name: "ID"},
			{Name: "Arch"},
			{Name: "State"},
			{Name: "Commit"},
			{Name: "Name"},
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
		omat, err := loadOmatConfig()
		if err != nil {
			util.Fatal(err)
		}

		details, err := awsutil.FindAndAssumeAdminRole(omat.BuildAccountSlug, omat)
		if err != nil {
			util.Fatal(err)
		}

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
