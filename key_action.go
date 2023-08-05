package main

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

type KeyAction struct {
	Key         *tcell.EventKey
	Description string
	Action      func()
}

func (k KeyAction) String() string {
	if k.Key.Key() == tcell.KeyRune {
		split := strings.Split(k.Key.Name(), "+")
		keyName := string(k.Key.Rune())
		parts := append(split[0:len(split)-1], keyName)
		return strings.Join(parts, "+")
	} else {
		return k.Key.Name()
	}
}
