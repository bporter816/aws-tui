package main

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/chroma/quick"
	"github.com/rivo/tview"
)

type Text struct {
	*tview.TextView
	highlightSyntax bool
	lang            string
}

func NewText(highlightSyntax bool, lang string) *Text {
	tv := tview.NewTextView()
	tv.SetDynamicColors(highlightSyntax)
	t := &Text{
		TextView:        tv,
		highlightSyntax: highlightSyntax,
		lang:            lang,
	}
	return t
}

func (t *Text) SetText(data string) {
	if data == "" {
		t.TextView.SetText("<empty>")
		return
	}

	if t.lang == "json" {
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(data), "", "  ")
		if err == nil {
			data = buf.String()
		}
	}

	if t.highlightSyntax {
		var buf bytes.Buffer
		err := quick.Highlight(&buf, data, t.lang, "terminal256", "solarized-dark256")
		if err == nil {
			data = buf.String()
			data = tview.TranslateANSI(data)
		}
	}

	t.TextView.SetText(data)
}
