package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"testing"
)

func TestAutoCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "TEST_CASE-1",
			expected: "Test case 1",
		},
		{
			input:    "HTTPS_ONLY",
			expected: "HTTPS only",
		},
		{
			input:    "hello there",
			expected: "Hello there",
		},
		{
			input:    "ONE_TWO_THREE_HTTP",
			expected: "One two three HTTP",
		},
	}

	for _, tc := range tests {
		got := AutoCase(tc.input)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}

func TestGetResourceNameFromArn(t *testing.T) {
	// examples from https://docs.aws.amazon.com/IAM/latest/UserGuide/reference-arns.html
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "arn:partition:service:region:account-id:resource-id",
			expected: "resource-id",
		},
		{
			input:    "arn:partition:service:region:account-id:resource-type/resource-id",
			expected: "resource-id",
		},
		{
			input:    "arn:partition:service:region:account-id:resource-type:resource-id",
			expected: "resource-id",
		},
		{
			input:    "arn:aws:dynamodb:us-east-2:123456789012:table/myDynamoDBTable",
			expected: "myDynamoDBTable",
		},
		{
			input:    "arn:aws:cloudfront::123456789012:distribution/myCloudfrontDistribution",
			expected: "myCloudfrontDistribution",
		},
	}

	for _, tc := range tests {
		arn, err := arn.Parse(tc.input)
		if err != nil {
			t.Fatalf("failed parsing arn: %v", tc.input)
		}
		got := GetResourceNameFromArn(arn)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
