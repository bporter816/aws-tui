package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type LambdaFunctions struct {
	*ui.Table
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

func (l LambdaFunctions) GetService() string {
	return "Lambda"
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
		tagsView := NewLambdaTags(l.repo, *arn, l.app)
		l.app.AddAndSwitch(tagsView)
	}
}

func (l LambdaFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		var name, status, runtime, codeSize, memory, timeout, layers, description string
		if v.FunctionName != nil {
			name = *v.FunctionName
		}
		status = utils.TitleCase(string(v.State))
		runtime = string(v.Runtime)
		codeSize = utils.FormatSize(v.CodeSize, 1)
		if v.MemorySize != nil {
			memory = utils.FormatSize(int64(*v.MemorySize)<<20, 1)
		}
		if v.Timeout != nil {
			timeout = strconv.FormatInt(int64(*v.Timeout), 10) + " sec"
		}
		layers = strconv.Itoa(len(v.Layers))
		if v.Description != nil {
			description = *v.Description
		}
		data = append(data, []string{
			name,
			status,
			runtime,
			codeSize,
			memory,
			timeout,
			layers,
			description,
		})
	}
	l.SetData(data)
}
