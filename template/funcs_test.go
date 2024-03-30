package template

import (
	"reflect"
	"testing"
)

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
