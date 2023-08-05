package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestKeyActionNotRune(t *testing.T) {
	k := KeyAction{
		Key:         tcell.NewEventKey(tcell.KeyEnter, 'a', tcell.ModNone),
		Description: "test",
		Action:      func() {},
	}
	assert.Equal(t, tcell.KeyNames[tcell.KeyEnter], k.String())
}

func TestKeyActionRuneWithNoModifiers(t *testing.T) {
	k := KeyAction{
		Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
		Description: "test",
		Action:      func() {},
	}
	assert.Equal(t, "a", k.String())
}

func TestKeyActionRuneWithCtrl(t *testing.T) {
	k := KeyAction{
		Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModCtrl),
		Description: "test",
		Action:      func() {},
	}
	assert.Equal(t, "Ctrl+a", k.String())
}
