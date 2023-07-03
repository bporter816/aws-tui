package main

import (
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	*tview.Table
	headers []string
}

func NewTable(headers []string, fixedRows, fixedCols int) *Table {
	tt := tview.NewTable()
	tt.SetFixed(fixedRows, fixedCols)
	tt.SetSelectable(true, false)
	tt.Select(1, 0)
	tt.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// TODO handle ctrl b
		if event.Key() == tcell.KeyHome || (event.Key() == tcell.KeyRune && event.Rune() == 'g') {
			tt.Select(1, 0)
			return nil
		}
		return event
	})
	for i, v := range headers {
		tt.SetCell(0, i, tview.NewTableCell(v))
	}
	t := &Table{
		Table:   tt,
		headers: headers,
	}
	return t
}

func (t *Table) SetData(data [][]string) {
	// TODO handle changing data size
	for r, v := range data {
		for c, vv := range v {
			t.SetCell(r+1, c, tview.NewTableCell(vv))
		}
	}
}

func (t Table) GetColSelection(col string) (string, error) {
	for i := 0; i < len(t.headers); i++ {
		if t.headers[i] == col {
			r, _ := t.GetSelection()
			return t.GetCell(r, i).Text, nil
		}
	}
	return "", errors.New("column not found")
}
