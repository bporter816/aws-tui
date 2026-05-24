package internal

import (
	"github.com/bporter816/aws-tui/internal/ui"
)

type TableView[T any] struct {
	*ui.Table
	service    string
	labels     []string
	items      []T
	load       func() ([]T, error)
	toRow      func(T) []string
	keyActions []KeyAction
}

func NewTableView[T any](
	service string,
	labels []string,
	columns []string,
	load func() ([]T, error),
	toRow func(T) []string,
	keyActions func(*TableView[T]) []KeyAction,
) *TableView[T] {
	v := &TableView[T]{
		Table:   ui.NewTable(columns, 1, 0),
		service: service,
		labels:  labels,
		load:    load,
		toRow:   toRow,
	}
	if keyActions != nil {
		v.keyActions = keyActions(v)
	}
	return v
}

func (t *TableView[T]) GetService() string        { return t.service }
func (t *TableView[T]) GetLabels() []string       { return t.labels }
func (t *TableView[T]) GetKeyActions() []KeyAction { return t.keyActions }

func (t *TableView[T]) Selected() (T, error) {
	row, err := t.GetRowSelection()
	if err != nil {
		var zero T
		return zero, err
	}
	return t.items[row-1], nil
}

func (t *TableView[T]) Render() error {
	items, err := t.load()
	if err != nil {
		return err
	}
	t.items = items
	rows := make([][]string, len(items))
	for i, item := range items {
		rows[i] = t.toRow(item)
	}
	t.SetData(rows)
	return nil
}
