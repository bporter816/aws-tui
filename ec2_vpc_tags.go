package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type EC2VPCTags struct {
	*ui.Table
	repo  *repo.EC2
	vpcId string
	app   *Application
}

func NewEC2VPCTags(repo *repo.EC2, vpcId string, app *Application) *EC2VPCTags {
	e := &EC2VPCTags{
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

func (e EC2VPCTags) GetService() string {
	return "EC2"
}

func (e EC2VPCTags) GetLabels() []string {
	return []string{e.vpcId, "Tags"}
}

func (e EC2VPCTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2VPCTags) Render() {
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
