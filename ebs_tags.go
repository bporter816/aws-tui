package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EBSTags struct {
	*ui.Table
	view.EBS
	resourceId string
	repo       *repo.EC2
	app        *Application
}

func NewEBSTags(resourceId string, repo *repo.EC2, app *Application) *EBSTags {
	e := &EBSTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		resourceId: resourceId,
		repo:       repo,
		app:        app,
	}
	return e
}

func (e EBSTags) GetLabels() []string {
	return []string{e.resourceId, "Tags"}
}

func (e EBSTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EBSTags) Render() {
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
