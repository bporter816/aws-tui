package internal

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestKeyActionString(t *testing.T) {
	tests := []struct {
		key      KeyAction
		expected string
	}{
		{
			key: KeyAction{
				Key:         tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone),
				Description: "test",
				Action:      func() {},
			},
			expected: tcell.KeyNames[tcell.KeyEnter],
		},
		{
			key: KeyAction{
				Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
				Description: "test",
				Action:      func() {},
			},
			expected: "a",
		},
		{
			key: KeyAction{
				Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModCtrl),
				Description: "test",
				Action:      func() {},
			},
			expected: "Ctrl+a",
		},
	}

	for _, tc := range tests {
		got := tc.key.String()
		if got != tc.expected {
			t.Fatalf("expected: %v, got: %v", tc.expected, got)
		}
	}
}
