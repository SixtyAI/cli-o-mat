package cmd

import (
	"fmt"
	"time"

	"github.com/FasterBetter/cli-o-mat/awsutil"
	"github.com/FasterBetter/cli-o-mat/util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

var (
	errCouldntLaunchInstance = errors.New("unable to launch instance")
)

// nolint: gochecknoglobals,gomnd
var launchCmd = &cobra.Command{
	Use:   "launch template-name keypair-name [subnet-id]",
	Short: "Launch an EC2 instance from a launch template.",
	Long: `Launch an EC2 instance from a launch template.

If you don't specify a subnet-id, the default subnet from the launch template will be used.`,
	Args: cobra.RangeArgs(2, 3),
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
		name := args[0]
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

// nolint: gochecknoglobals
var (
	launchVersion string
	launchType    string
)

// nolint: gochecknoinits
func init() {
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&launchVersion, "version", "", "", "Version of launch template to use (default: $Latest)")
	launchCmd.Flags().StringVarP(&launchType, "type", "", "", "Instance type to launch (default from launch template)")
}
