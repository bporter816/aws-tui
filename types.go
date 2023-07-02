package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type KeyAction struct {
	Key         *tcell.EventKey
	Description string
	Action      func()
}

type Component interface {
	tview.Primitive
	GetName() string // TODO maybe []string?
	GetKeyActions() []KeyAction
	Render()
}
