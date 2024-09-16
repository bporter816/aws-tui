package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type LambdaFunctions struct {
	*ui.Table
	view.Lambda
	repo  *repo.Lambda
	app   *Application
	model []model.LambdaFunction
}

func NewLambdaFunctions(repo *repo.Lambda, app *Application) *LambdaFunctions {
	l := &LambdaFunctions{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"RUNTIME",
			"CODE SIZE",
			"MEMORY",
			"TIMEOUT",
			"LAYERS",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return l
}

func (l LambdaFunctions) GetLabels() []string {
	return []string{"Functions"}
}

func (l LambdaFunctions) tagsHandler() {
	row, err := l.GetRowSelection()
	if err != nil {
		return
	}
	if arn := l.model[row-1].FunctionArn; arn != nil {
		tagsView := NewTags(l.repo, l.GetService(), *arn, l.app)
		l.app.AddAndSwitch(tagsView)
	}
}

func (l LambdaFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      l.tagsHandler,
		},
	}
}

func (l *LambdaFunctions) Render() {
	model, err := l.repo.ListFunctions()
	if err != nil {
		panic(err)
	}
	l.model = model

	var data [][]string
	for _, v := range model {
		var memory, timeout string
		if v.MemorySize != nil {
			memory = utils.FormatSize(int64(*v.MemorySize)<<20, 1)
		}
		if v.Timeout != nil {
			timeout = strconv.FormatInt(int64(*v.Timeout), 10) + " sec"
		}
		data = append(data, []string{
			utils.DerefString(v.FunctionName, ""),
			utils.TitleCase(string(v.State)),
			string(v.Runtime),
			utils.FormatSize(v.CodeSize, 1),
			memory,
			timeout,
			strconv.Itoa(len(v.Layers)),
			utils.DerefString(v.Description, ""),
		})
	}
	l.SetData(data)
}
