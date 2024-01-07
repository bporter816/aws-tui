package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type LambdaTags struct {
	*ui.Table
	view.Lambda
	repo *repo.Lambda
	id   string
	app  *Application
}

func NewLambdaTags(repo *repo.Lambda, id string, app *Application) *LambdaTags {
	l := &LambdaTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		id:   id,
		app:  app,
	}
	return l
}

func (l LambdaTags) GetLabels() []string {
	arn, err := arn.Parse(l.id)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (l LambdaTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (l LambdaTags) Render() {
	model, err := l.repo.ListTags(l.id)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	l.SetData(data)
}
