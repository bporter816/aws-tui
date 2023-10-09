package utils

import (
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
	}

	for _, tc := range tests {
		got := AutoCase(tc.input)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
