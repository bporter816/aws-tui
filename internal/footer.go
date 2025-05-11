package internal

import (
	"github.com/rivo/tview"
	"strings"
)

type Footer struct {
	*tview.TextView
	app *Application
}

func NewFooter(app *Application) *Footer {
	f := &Footer{
		TextView: tview.NewTextView().SetDynamicColors(true),
		app:      app,
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
