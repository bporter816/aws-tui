package main

import (
	"github.com/rivo/tview"
)

type Component interface {
	tview.Primitive
	GetService() string
	GetLabels() []string
	GetKeyActions() []KeyAction
	Render()
}
