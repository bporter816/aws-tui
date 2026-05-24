package internal

import (
	"context"

	"github.com/rivo/tview"
)

type Component interface {
	tview.Primitive
	GetService() string
	GetLabels() []string
	GetKeyActions() []KeyAction
	Render(context.Context) error
}
