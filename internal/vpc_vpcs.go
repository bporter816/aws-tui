package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type VPCVPCs struct {
	*ui.Table
	view.VPC
	repo *repo.EC2
	app  *Application
}

func NewVPCVPCs(repo *repo.EC2, app *Application) *VPCVPCs {
	e := &VPCVPCs{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"STATE",
			"IPV4 CIDR",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e VPCVPCs) GetLabels() []string {
	return []string{"VPCs"}
}

func (e VPCVPCs) tagsHandler() {
	vpcId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), vpcId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e VPCVPCs) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e VPCVPCs) Render() {
	model, err := e.repo.ListVPCs()
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
			utils.DerefString(v.VpcId, ""),
			utils.TitleCase(string(v.State)),
			utils.DerefString(v.CidrBlock, ""),
		})
	}
	e.SetData(data)
}
