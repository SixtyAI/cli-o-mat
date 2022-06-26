package awsutil

import "github.com/aws/aws-sdk-go/aws"

func DefaultToString(value *bool) string {
	var isDefault string
	if aws.BoolValue(value) {
		isDefault = "yes"
	} else {
		isDefault = "" // Blank so the `yes` value stands out more.
	}

	return isDefault
}
