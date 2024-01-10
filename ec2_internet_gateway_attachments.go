package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type EC2InternetGatewayAttachments struct {
	*ui.Table
	view.EC2
	repo              *repo.EC2
	internetGatewayId string
	app               *Application
}

func NewEC2InternetGatewayAttachments(repo *repo.EC2, internetGatewayId string, app *Application) *EC2InternetGatewayAttachments {
	e := &EC2InternetGatewayAttachments{
		Table: ui.NewTable([]string{
			"VPC ID",
			"STATE",
		}, 1, 0),
		repo:              repo,
		internetGatewayId: internetGatewayId,
		app:               app,
	}
	return e
}

func (e EC2InternetGatewayAttachments) GetLabels() []string {
	return []string{e.internetGatewayId, "Attachments"}
}

func (e EC2InternetGatewayAttachments) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2InternetGatewayAttachments) Render() {
	model, err := e.repo.ListInternetGatewayAttachments(e.internetGatewayId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var vpcId, state string
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		state = utils.AutoCase(string(v.State))
		data = append(data, []string{
			vpcId,
			state,
		})
	}
	e.SetData(data)
}
