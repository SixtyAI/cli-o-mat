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

func showImages(images []*ec2.Image) {
	// N.B. It's a bit rude to modify our parameter in-place in a function the caller expects to just
	// show things, but since we're not reusing the data structure yet, I'm not gonna add the overhead
	// of the slice copy.
	sort.Slice(images, func(i, j int) bool {
		return aws.StringValue(images[i].Name) < aws.StringValue(images[j].Name)
	})

	const imagesFormat = "%-64s %-6s %-21s %s\n"

	fmt.Printf(imagesFormat, "Name", "Arch", "ID", "State")

	for _, image := range images {
		fmt.Printf(imagesFormat, aws.StringValue(image.Name), aws.StringValue(image.Architecture),
			aws.StringValue(image.ImageId), aws.StringValue(image.State))
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

		showImages(imageList.Images)
	},
}

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(imagesCmd)
}
