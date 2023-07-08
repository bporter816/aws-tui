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
	GetService() string
	GetLabels() []string
	GetKeyActions() []KeyAction
	Render()
}
