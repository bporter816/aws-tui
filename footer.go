package main

import (
	"github.com/bporter816/aws-tui/ui"
	"strings"
)

type Footer struct {
	*ui.Text
	app *Application
}

func NewFooter(app *Application) *Footer {
	f := &Footer{
		Text: ui.NewText(false, ""),
		app:  app,
	}
	return f
}

func (f Footer) Render() {
	var names []string
	for i, v := range f.app.components {
		if i < 2 {
			names = append(names, v.GetService())
		}
		names = append(names, v.GetLabels()...)
	}
	str := strings.Join(names, " > ")
	f.SetText(str)
}
