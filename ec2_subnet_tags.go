package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2SubnetTags struct {
	*ui.Table
	view.EC2
	repo     *repo.EC2
	subnetId string
	app      *Application
}

func NewEC2SubnetTags(repo *repo.EC2, subnetId string, app *Application) *EC2SubnetTags {
	e := &EC2SubnetTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:     repo,
		subnetId: subnetId,
		app:      app,
	}
	return e
}

func (e EC2SubnetTags) GetLabels() []string {
	return []string{e.subnetId, "Tags"}
}

func (e EC2SubnetTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2SubnetTags) Render() {
	model, err := e.repo.ListSubnetTags(e.subnetId)
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
