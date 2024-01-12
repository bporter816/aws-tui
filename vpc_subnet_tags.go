package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type VPCSubnetTags struct {
	*ui.Table
	view.VPC
	repo     *repo.EC2
	subnetId string
	app      *Application
}

func NewVPCSubnetTags(repo *repo.EC2, subnetId string, app *Application) *VPCSubnetTags {
	e := &VPCSubnetTags{
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

func (e VPCSubnetTags) GetLabels() []string {
	return []string{e.subnetId, "Tags"}
}

func (e VPCSubnetTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e VPCSubnetTags) Render() {
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
