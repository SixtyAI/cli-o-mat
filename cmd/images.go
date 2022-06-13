package cmd

import (
	"fmt"
	"sort"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
)

func showImages(shortHashes bool, images []*ec2.Image) {
	// N.B. It's a bit rude to modify our parameter in-place in a function the caller expects to just
	// show things, but since we're not reusing the data structure yet, I'm not gonna add the overhead
	// of the slice copy.
	sort.Slice(images, func(i, j int) bool {
		return aws.StringValue(images[i].Name) < aws.StringValue(images[j].Name)
	})

	hashLen := 40
	if shortHashes {
		hashLen = 7
	}

	imagesFormat := fmt.Sprintf("%%-21s %%-6s %%-9s %%-%ds %%s\n", hashLen)

	fmt.Printf(imagesFormat, "ID", "Arch", "State", "Commit", "Name")

	for _, image := range images {
		var commit string

		for _, value := range image.Tags {
			if aws.StringValue(value.Key) == "BuildCommit" {
				commit = aws.StringValue(value.Value)

				if shortHashes {
					commit = commit[0:7]
				}

				break
			}
		}

		fmt.Printf(imagesFormat, aws.StringValue(image.ImageId), aws.StringValue(image.Architecture),
			aws.StringValue(image.State), commit, aws.StringValue(image.Name))
	}
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

		fmt.Printf("\n\n\nGOT TAGS: %+v\n\n\n", imageList.Images[0].Tags)

		showImages(shortHashes, imageList.Images)
	},
}

// nolint: gochecknoglobals
var shortHashes bool

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(imagesCmd)

	imagesCmd.Flags().BoolVarP(&shortHashes, "short", "", false, "Shorten git commit SHAs")
}
