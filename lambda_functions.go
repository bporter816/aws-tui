package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"strconv"
)

type LambdaFunctions struct {
	*ui.Table
	repo *repo.Lambda
	app  *Application
}

func NewLambdaFunctions(repo *repo.Lambda, app *Application) *LambdaFunctions {
	l := &LambdaFunctions{
		Table: ui.NewTable([]string{
			"NAME",
			"RUNTIME",
			"LAYERS",
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

func (l LambdaFunctions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (l LambdaFunctions) Render() {
	model, err := l.repo.ListFunctions()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, runtime, layers string
		if v.FunctionName != nil {
			name = *v.FunctionName
		}
		runtime = string(v.Runtime)
		layers = strconv.Itoa(len(v.Layers))
		data = append(data, []string{
			name,
			runtime,
			layers,
		})
	}
	l.SetData(data)
}
