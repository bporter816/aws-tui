package utils

import (
	"fmt"
	"testing"

	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func TestJoinKMSGrantOperations(t *testing.T) {
	tests := []struct {
		operations []kmsTypes.GrantOperation
		sep        string
		expected   string
	}{
		{
			operations: []kmsTypes.GrantOperation{
				kmsTypes.GrantOperationEncrypt,
				kmsTypes.GrantOperationDecrypt,
			},
			sep:      ", ",
			expected: fmt.Sprintf("%v, %v", kmsTypes.GrantOperationEncrypt, kmsTypes.GrantOperationDecrypt),
		},
	}

	for _, tc := range tests {
		got := JoinKMSGrantOperations(tc.operations, tc.sep)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}

func TestFormatKMSAliases(t *testing.T) {
	tests := []struct {
		aliases  []string
		expected string
	}{
		{
			aliases:  []string{"a"},
			expected: "a",
		},
		{
			aliases:  []string{"a", "b"},
			expected: "a + 1 more",
		},
		{
			aliases:  []string{"a", "b", "c"},
			expected: "a + 2 more",
		},
	}

	for _, tc := range tests {
		got := FormatKMSAliases(tc.aliases)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
