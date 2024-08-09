package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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
		var name, id, state, ipv4CIDR string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		}
		if v.VpcId != nil {
			id = *v.VpcId
		}
		state = utils.TitleCase(string(v.State))
		if v.CidrBlock != nil {
			ipv4CIDR = *v.CidrBlock
		}
		data = append(data, []string{
			name,
			id,
			state,
			ipv4CIDR,
		})
	}
	e.SetData(data)
}
