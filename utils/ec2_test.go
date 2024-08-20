package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"testing"
)

func TestLookupEC2Tag(t *testing.T) {
	tests := []struct {
		tags        []ec2Types.Tag
		key         string
		expectedVal string
		expectedOk  bool
	}{
		{
			tags: []ec2Types.Tag{
				{
					Key:   aws.String("k"),
					Value: aws.String("v"),
				},
			},
			key:         "k",
			expectedVal: "v",
			expectedOk:  true,
		},
		{
			tags: []ec2Types.Tag{
				{
					Key:   aws.String("k"),
					Value: aws.String("v"),
				},
			},
			key:         "k2",
			expectedVal: "",
			expectedOk:  false,
		},
		{
			tags: []ec2Types.Tag{
				{
					Key:   aws.String("k"),
					Value: aws.String("v"),
				},
				{
					Key:   aws.String("k2"),
					Value: aws.String("v2"),
				},
			},
			key:         "k2",
			expectedVal: "v2",
			expectedOk:  true,
		},
	}

	for _, tc := range tests {
		val, ok := LookupEC2Tag(tc.tags, tc.key)
		if val != tc.expectedVal || ok != tc.expectedOk {
			t.Fatalf("expected: (%v, %v), got: (%v, %v)", tc.expectedVal, tc.expectedOk, val, ok)
		}
	}
}
