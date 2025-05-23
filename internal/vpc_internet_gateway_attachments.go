package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type VPCInternetGatewayAttachments struct {
	*ui.Table
	view.VPC
	repo              *repo.EC2
	internetGatewayId string
	app               *Application
}

func NewVPCInternetGatewayAttachments(repo *repo.EC2, internetGatewayId string, app *Application) *VPCInternetGatewayAttachments {
	e := &VPCInternetGatewayAttachments{
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

func (e VPCInternetGatewayAttachments) GetLabels() []string {
	return []string{e.internetGatewayId, "Attachments"}
}

func (e VPCInternetGatewayAttachments) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e VPCInternetGatewayAttachments) Render() {
	model, err := e.repo.ListInternetGatewayAttachments(e.internetGatewayId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.VpcId, ""),
			utils.AutoCase(string(v.State)),
		})
	}
	e.SetData(data)
}
