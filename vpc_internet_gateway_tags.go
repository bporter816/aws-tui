package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type VPCInternetGatewayTags struct {
	*ui.Table
	view.VPC
	repo              *repo.EC2
	internetGatewayId string
	app               *Application
}

func NewVPCInternetGatewayTags(repo *repo.EC2, internetGatewayId string, app *Application) *VPCInternetGatewayTags {
	e := &VPCInternetGatewayTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:              repo,
		internetGatewayId: internetGatewayId,
		app:               app,
	}
	return e
}

func (e VPCInternetGatewayTags) GetLabels() []string {
	return []string{e.internetGatewayId, "Tags"}
}

func (e VPCInternetGatewayTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e VPCInternetGatewayTags) Render() {
	model, err := e.repo.ListInternetGatewayTags(e.internetGatewayId)
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
