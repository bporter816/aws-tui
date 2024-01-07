package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2InstanceTags struct {
	*ui.Table
	view.EC2
	repo       *repo.EC2
	instanceId string
	app        *Application
}

func NewEC2InstanceTags(repo *repo.EC2, instanceId string, app *Application) *EC2InstanceTags {
	e := &EC2InstanceTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:       repo,
		instanceId: instanceId,
		app:        app,
	}
	return e
}

func (e EC2InstanceTags) GetLabels() []string {
	return []string{e.instanceId, "Tags"}
}

func (e EC2InstanceTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2InstanceTags) Render() {
	model, err := e.repo.ListInstanceTags(e.instanceId)
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
