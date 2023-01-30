package main

import (
	"github.com/rivo/tview"
)

type TableRenderable interface {
	GetHeaders() []string
	Render() ([][]string, error)
}

type TableView struct {
	table   *tview.Table
	name    string
	content TableRenderable
}

func NewTableView(name string, content TableRenderable) *TableView {
	t := &TableView{
		table:   tview.NewTable(),
		name:    name,
		content: content,
	}
	t.table.SetFixed(1, 0)
	t.table.SetSelectable(true, false)
	t.table.SetBorder(true)
	t.table.Box.SetTitle(" " + name + " ")
	for i, v := range content.GetHeaders() {
		t.table.SetCell(0, i, tview.NewTableCell(v))
	}
	t.Update()
	return t
}

func (t TableView) GetTable() *tview.Table {
	return t.table
}

func (t TableView) Update() {
	data, err := t.content.Render()
	if err != nil {
		panic(err)
	}
	for r, v := range data {
		for c, vv := range v {
			t.table.SetCell(r+1, c, tview.NewTableCell(vv))
		}
	}
}
