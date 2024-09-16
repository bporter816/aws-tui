package utils

import (
	"testing"
)

func TestDerefString(t *testing.T) {
	str := "test"
	tests := []struct {
		v        *string
		d        string
		expected string
	}{
		{
			v:        &str,
			d:        "default",
			expected: str,
		},
		{
			v:        nil,
			d:        "default",
			expected: "default",
		},
	}

	for _, tc := range tests {
		got := DerefString(tc.v, tc.d)
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
