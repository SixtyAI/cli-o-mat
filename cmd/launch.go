package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"
)

var errCouldntLaunchInstance = errors.New("unable to launch instance")

//nolint: gochecknoglobals,gomnd
var launchCmd = &cobra.Command{
	Use:   "launch template-name keypair-name [subnet-id]",
	Short: "Launch an EC2 instance from a launch template.",
	Long: `Launch an EC2 instance from a launch template.

If you don't specify a subnet-id, the default subnet from the launch template will be used.`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		omat := loadOmatConfig()

		details := awsutil.FindAndAssumeAdminRole(omat.DeployAccountSlug, omat)

		ec2Client := ec2.New(details.Session, details.Config)
		namePrefix := args[0]
		keypair := args[1]

		var subnetID *string
		if len(args) == 3 {
			subnetID = aws.String(args[2])
		}

		if launchVersion == "" {
			launchVersion = "$Latest"
		}

		var instanceType *string
		if launchType != "" {
			instanceType = aws.String(launchType)
		}

		templates, err := awsutil.FetchLaunchTemplates(ec2Client, nil)
		if err != nil {
			util.Fatal(err)
		}
		candidates := make([]string, 0)
		for _, template := range templates {
			templateName := aws.StringValue(template.LaunchTemplateName)
			if strings.HasPrefix(templateName, namePrefix) {
				candidates = append(candidates, templateName)
			}
		}

		if len(candidates) == 0 {
			fmt.Printf("Found the following launch templates, none of which match specified prefix:\n")
			for _, template := range templates {
				fmt.Printf("\t%s\n", aws.StringValue(template.LaunchTemplateName))
			}
			util.Fatal(errors.New("no matching launch templates found"))
		} else if len(candidates) > 1 {
			fmt.Printf("Found the following launch templates matching specified prefix:\n")
			for _, candidate := range candidates {
				fmt.Printf("\t%s\n", candidate)
			}
			util.Fatal(errors.New("multiple launch templates found"))
		}
		name := candidates[0]
		fmt.Printf("Using launch template %s...\n", name)

		resp, err := ec2Client.RunInstances(&ec2.RunInstancesInput{
			LaunchTemplate: &ec2.LaunchTemplateSpecification{
				LaunchTemplateName: &name,
				Version:            aws.String(launchVersion),
			},
			InstanceType:                      instanceType,
			InstanceInitiatedShutdownBehavior: aws.String("terminate"),
			KeyName:                           aws.String(keypair),

			MinCount: aws.Int64(1),
			MaxCount: aws.Int64(1),
			SubnetId: subnetID,
		})
		if err != nil {
			util.Fatal(err)
		}

		if len(resp.Instances) != 1 {
			util.Fatal(errCouldntLaunchInstance)
		}

		fmt.Printf("Launching instance %s...\n", aws.StringValue(resp.Instances[0].InstanceId))
		fmt.Printf("Waiting for instance to public IP...\n")

		counter := 0
		instanceIds := []*string{resp.Instances[0].InstanceId}
		publicIP := ""

		for {
			<-time.After(1 * time.Second)

			counter++
			if counter > 30 {
				break
			}

			resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: instanceIds,
			})
			if err != nil {
				util.Fatal(err)
			}

			if aws.StringValue(resp.Reservations[0].Instances[0].PublicIpAddress) != "" {
				publicIP = aws.StringValue(resp.Reservations[0].Instances[0].PublicIpAddress)

				break
			}
		}

		if publicIP != "" {
			fmt.Printf("Public IP: %s\n", publicIP)
		} else {
			fmt.Printf("Couldn't determine public IP.\n")
		}
	},
}

//nolint: gochecknoglobals
var (
	launchVersion string
	launchType    string
)

//nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&launchVersion, "version", "", "", "Version of launch template to use (default: $Latest)")
	launchCmd.Flags().StringVarP(&launchType, "type", "", "", "Instance type to launch (default from launch template)")
}
