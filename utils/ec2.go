package utils

import (
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TODO add tests
func LookupEC2Tag(tags []ec2Types.Tag, key string) (string, bool) {
	for _, v := range tags {
		if *v.Key == key {
			return *v.Value, true
		}
	}
	return "", false
}
