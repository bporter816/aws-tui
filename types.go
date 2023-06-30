package main

import (
	"github.com/gdamore/tcell/v2"
)

type KeyAction struct {
	Key         *tcell.EventKey
	Description string
	Action      func()
}

type Component interface {
	GetName() string // TODO maybe []string?
	GetKeyActions() []KeyAction
	Render()
}
