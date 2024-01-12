package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2Tags struct {
	*ui.Table
	view.EC2
	repo       *repo.EC2
	resourceId string
	app        *Application
}

func NewEC2Tags(repo *repo.EC2, resourceId string, app *Application) *EC2Tags {
	e := &EC2Tags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:       repo,
		resourceId: resourceId,
		app:        app,
	}
	return e
}

func (e EC2Tags) GetLabels() []string {
	return []string{e.resourceId, "Tags"}
}

func (e EC2Tags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2Tags) Render() {
	model, err := e.repo.ListTags(e.resourceId)
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
