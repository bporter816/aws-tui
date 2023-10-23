package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"reflect"
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

func TestSimplifyFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{
			input:    1.0,
			expected: "1",
		},
		{
			input:    1.5,
			expected: "1.5",
		},
	}

	for _, tc := range tests {
		got := SimplifyFloat(tc.input)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}

func TestBoolToString(t *testing.T) {
	tests := []struct {
		boolean  bool
		yesStr   string
		noStr    string
		expected string
	}{
		{
			boolean:  true,
			yesStr:   "yes",
			noStr:    "no",
			expected: "yes",
		},
		{
			boolean:  false,
			yesStr:   "yes",
			noStr:    "no",
			expected: "no",
		},
	}

	for _, tc := range tests {
		got := BoolToString(tc.boolean, tc.yesStr, tc.noStr)
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
			input:    "arn:aws:cloudfront::123456789012:distribution/myCloudFrontDistribution",
			expected: "myCloudFrontDistribution",
		},
		{
			input:    "arn:aws:s3:::bucket/path/to/object",
			expected: "/path/to/object",
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

func TestFormatSerial(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte{},
			expected: "",
		},
		{
			input:    []byte{0x00},
			expected: "00",
		},
		{
			input:    []byte{0x12, 0x34},
			expected: "12:34",
		},
	}

	for _, tc := range tests {
		got := FormatSerial(tc.input)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		str      string
		length   int
		expected []string
	}{
		{
			str:      "",
			length:   2,
			expected: []string{},
		},
		{
			str:      "aabbccdd",
			length:   2,
			expected: []string{"aa", "bb", "cc", "dd"},
		},
		{
			str:      "aabbccd",
			length:   2,
			expected: []string{"aa", "bb", "cc", "d"},
		},
		{
			str:      "aaa",
			length:   5,
			expected: []string{"aaa"},
		},
	}

	for _, tc := range tests {
		got := Chunk(tc.str, tc.length)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size      int64
		precision int
		expected  string
	}{
		{
			size:      0,
			precision: 1,
			expected:  "0 B",
		},
		{
			size:      999,
			precision: 1,
			expected:  "999 B",
		},
		{
			size:      1024,
			precision: 1,
			expected:  "1.0 KiB",
		},
		{
			size:      1572864,
			precision: 1,
			expected:  "1.5 MiB",
		},
	}

	for _, tc := range tests {
		got := FormatSize(tc.size, tc.precision)
		if got != tc.expected {
			t.Fatalf("expected %v, got: %v", tc.expected, got)
		}
	}
}
