package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2InternetGatewayTags struct {
	*ui.Table
	view.EC2
	repo              *repo.EC2
	internetGatewayId string
	app               *Application
}

func NewEC2InternetGatewayTags(repo *repo.EC2, internetGatewayId string, app *Application) *EC2InternetGatewayTags {
	e := &EC2InternetGatewayTags{
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

func (e EC2InternetGatewayTags) GetLabels() []string {
	return []string{e.internetGatewayId, "Tags"}
}

func (e EC2InternetGatewayTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2InternetGatewayTags) Render() {
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
