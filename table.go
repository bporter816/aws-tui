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
	r, _ := t.GetSelection()
	if r == 0 {
		return "", errors.New("cannot select row 0")
	}
	if r >= t.GetRowCount() {
		return "", errors.New("selection out of range")
	}
	for i := 0; i < len(t.headers); i++ {
		if t.headers[i] == col {
			return t.GetCell(r, i).Text, nil
		}
	}
	return "", errors.New("column not found")
}

func (t Table) GetRowSelection() (int, error) {
	r, _ := t.GetSelection()
	if r == 0 {
		return 0, errors.New("cannot select row 0")
	}
	if r >= t.GetRowCount() {
		return 0, errors.New("selection out of range")
	}
	return r, nil
}
