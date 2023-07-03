package main

import (
	"bytes"
	"github.com/rivo/tview"
	"github.com/alecthomas/chroma/quick"
)

type Text struct {
	*tview.TextView
	highlightSyntax bool
	highlightLang   string
}

func NewText(highlightSyntax bool, highlightLang string) *Text {
	tv := tview.NewTextView()
	tv.SetDynamicColors(highlightSyntax)
	t := &Text{
		TextView:        tv,
		highlightSyntax: highlightSyntax,
		highlightLang:   highlightLang,
	}
	return t
}

func (t *Text) SetText(data string) {
	// TODO handle [ ] characters in input
	if t.highlightSyntax {
		var buf bytes.Buffer
		err := quick.Highlight(&buf, data, t.highlightLang, "terminal256", "solarized-dark256")
		if err != nil {
			panic(err)
		}
		data = buf.String()
		data = tview.TranslateANSI(data)
	}
	t.TextView.SetText(data)
}
