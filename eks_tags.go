package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type EKSTags struct {
	*ui.Table
	view.EKS
	repo *repo.EKS
	arn  string
	app  *Application
}

func NewEKSTags(repo *repo.EKS, arn string, app *Application) *EKSTags {
	e := &EKSTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		arn:  arn,
		app:  app,
	}
	return e
}

func (e EKSTags) GetLabels() []string {
	arn, err := arn.Parse(e.arn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (e EKSTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EKSTags) Render() {
	model, err := e.repo.ListTags(e.arn)
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
	e.SetData(data)
}
