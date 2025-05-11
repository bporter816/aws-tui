package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type VPCInternetGateways struct {
	*ui.Table
	view.VPC
	repo *repo.EC2
	app  *Application
}

func NewVPCInternetGateways(repo *repo.EC2, app *Application) *VPCInternetGateways {
	e := &VPCInternetGateways{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"OWNER",
			"ATTACHMENTS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e VPCInternetGateways) GetLabels() []string {
	return []string{"Internet Gateways"}
}

func (e VPCInternetGateways) attachmentsHandler() {
	internetGatewayId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	attachmentsView := NewVPCInternetGatewayAttachments(e.repo, internetGatewayId, e.app)
	e.app.AddAndSwitch(attachmentsView)
}

func (e VPCInternetGateways) tagsHandler() {
	internetGatewayId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), internetGatewayId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e VPCInternetGateways) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Attachments",
			Action:      e.attachmentsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e VPCInternetGateways) Render() {
	model, err := e.repo.ListInternetGateways()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name string
		if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
			name = n
		}
		data = append(data, []string{
			name,
			utils.DerefString(v.InternetGatewayId, ""),
			utils.DerefString(v.OwnerId, ""),
			strconv.Itoa(len(v.Attachments)),
		})
	}
	e.SetData(data)
}
