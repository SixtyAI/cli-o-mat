package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
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
		omatConfig, err := loadOmatConfig()
		if err != nil {
			util.Fatal(err)
		}

		paramPrefix := omatConfig.Prefix()

		rootSess := session.Must(session.NewSession())
		ssmClient := ssm.New(rootSess, aws.NewConfig().WithRegion(omatConfig.Region))

		roleParamName := fmt.Sprintf("%s/ci-cd/roles/admin", paramPrefix)
		roleParam, err := ssmClient.GetParameter(&ssm.GetParameterInput{
			Name: aws.String(roleParamName),
		})
		if err != nil {
			if strings.HasPrefix(err.Error(), "ParameterNotFound") {
				util.Fatalf("Could not find role parameter: %s\n", roleParamName)
			} else {
				util.Fatal(err)
			}
		}

		arn := aws.StringValue(roleParam.Parameter.Value)
		if arn != "" {
			fmt.Printf("Using ARN: %s (from %s)\n", *roleParam.Parameter.Value, roleParamName)
		} else {
			util.Fatalf("SSM paramater '%s' was empty.\n", roleParamName)
		}

		assumedSess := session.Must(session.NewSession())
		creds := stscreds.NewCredentials(assumedSess, arn)
		ec2Client := ec2.New(assumedSess, &aws.Config{Credentials: creds})

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
