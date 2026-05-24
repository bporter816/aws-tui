package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewVPCInternetGateways(repo *repo.EC2, app *Application) *TableView[model.EC2InternetGateway] {
	return NewTableView(
		"VPC",
		[]string{"Internet Gateways"},
		[]string{"NAME", "ID", "OWNER", "ATTACHMENTS"},
		repo.ListInternetGateways,
		func(v model.EC2InternetGateway) []string {
			var name string
			if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = n
			}
			return []string{
				name,
				utils.DerefString(v.InternetGatewayId, ""),
				utils.DerefString(v.OwnerId, ""),
				strconv.Itoa(len(v.Attachments)),
			}
		},
		func(v *TableView[model.EC2InternetGateway]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Attachments",
					Action: func() {
						internetGatewayId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewVPCInternetGatewayAttachments(repo, internetGatewayId, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						internetGatewayId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), internetGatewayId, app))
					},
				},
			}
		},
	)
}

func init() {
	register("VPC.Internet Gateways", func(r Repos, app *Application) Component {
		return NewVPCInternetGateways(r.EC2, app)
	})
}
