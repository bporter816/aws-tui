package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type VPCTags struct {
	*ui.Table
	view.VPC
	repo       *repo.EC2
	resourceId string
	app        *Application
}

func NewVPCTags(repo *repo.EC2, resourceId string, app *Application) *VPCTags {
	e := &VPCTags{
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

func (e VPCTags) GetLabels() []string {
	return []string{e.resourceId, "Tags"}
}

func (e VPCTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e VPCTags) Render() {
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
