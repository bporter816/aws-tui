package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type VPCVPCTags struct {
	*ui.Table
	view.VPC
	repo  *repo.EC2
	vpcId string
	app   *Application
}

func NewVPCVPCTags(repo *repo.EC2, vpcId string, app *Application) *VPCVPCTags {
	e := &VPCVPCTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:  repo,
		vpcId: vpcId,
		app:   app,
	}
	return e
}

func (e VPCVPCTags) GetLabels() []string {
	return []string{e.vpcId, "Tags"}
}

func (e VPCVPCTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e VPCVPCTags) Render() {
	model, err := e.repo.ListVPCTags(e.vpcId)
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
