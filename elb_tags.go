package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type ELBTags struct {
	*ui.Table
	view.ELB
	repo         *repo.ELB
	resourceArn  string
	resourceName string
	app          *Application
}

func NewELBTags(repo *repo.ELB, resourceArn string, resourceName string, app *Application) *ELBTags {
	e := &ELBTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:         repo,
		resourceArn:  resourceArn,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ELBTags) GetLabels() []string {
	return []string{e.resourceName, "Tags"}
}

func (e ELBTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTags) Render() {
	model, err := e.repo.ListTags(e.resourceArn)
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
