package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewVPCVPCs(repo *repo.EC2, app *Application) *TableView[model.EC2VPC] {
	return NewTableView(
		"VPC",
		[]string{"VPCs"},
		[]string{"NAME", "ID", "STATE", "IPV4 CIDR"},
		repo.ListVPCs,
		func(v model.EC2VPC) []string {
			var name string
			if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = n
			}
			return []string{
				name,
				utils.DerefString(v.VpcId, ""),
				utils.TitleCase(string(v.State)),
				utils.DerefString(v.CidrBlock, ""),
			}
		},
		func(v *TableView[model.EC2VPC]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						vpcId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), vpcId, app))
					},
				},
			}
		},
	)
}
