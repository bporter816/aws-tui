package ui

import (
	"bytes"
	"encoding/json"
	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Text struct {
	*tview.TextView
	HighlightSyntax bool
	Lang            string
}

func NewText(highlightSyntax bool, lang string) *Text {
	tv := tview.NewTextView()
	tv.SetBackgroundColor(tcell.ColorDefault)
	tv.SetDynamicColors(highlightSyntax)
	t := &Text{
		TextView:        tv,
		HighlightSyntax: highlightSyntax,
		Lang:            lang,
	}
	return t
}

func (t *Text) SetText(data string) {
	if data == "" {
		t.TextView.SetText("<empty>")
		return
	}

	if t.Lang == "json" {
		var buf bytes.Buffer
		err := json.Indent(&buf, []byte(data), "", "  ")
		if err == nil {
			data = buf.String()
		}
	}

	if t.HighlightSyntax {
		var buf bytes.Buffer
		err := quick.Highlight(&buf, data, t.Lang, "terminal256", "solarized-dark256")
		if err == nil {
			data = buf.String()
			data = tview.TranslateANSI(data)
		}
	}

	t.TextView.SetText(data)
}
