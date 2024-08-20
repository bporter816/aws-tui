package utils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	gaTypes "github.com/aws/aws-sdk-go-v2/service/globalaccelerator/types"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		input    []gaTypes.PortRange
		expected string
	}{
		{
			input: []gaTypes.PortRange{
				{
					FromPort: aws.Int32(80),
					ToPort:   aws.Int32(80),
				},
				{
					FromPort: aws.Int32(443),
					ToPort:   aws.Int32(443),
				},
			},
			expected: "80, 443",
		},
		{
			input: []gaTypes.PortRange{
				{
					FromPort: aws.Int32(1000),
					ToPort:   aws.Int32(2000),
				},
			},
			expected: "1000-2000",
		},
		{
			input: []gaTypes.PortRange{
				{
					FromPort: aws.Int32(1),
					ToPort:   aws.Int32(2),
				},
				{
					FromPort: aws.Int32(3),
					ToPort:   aws.Int32(3),
				},
				{
					FromPort: aws.Int32(4),
					ToPort:   aws.Int32(5),
				},
			},
			expected: "1-2, 3, 4-5",
		},
	}

	for _, tc := range tests {
		got := FormatGlobalAcceleratorPortRanges(tc.input)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
