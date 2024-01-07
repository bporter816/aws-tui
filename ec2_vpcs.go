package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2VPCs struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2VPCs(repo *repo.EC2, app *Application) *EC2VPCs {
	e := &EC2VPCs{
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

func (e EC2VPCs) GetLabels() []string {
	return []string{"VPCs"}
}

func (e EC2VPCs) tagsHandler() {
	vpcId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2VPCTags(e.repo, vpcId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2VPCs) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2VPCs) Render() {
	model, err := e.repo.ListVPCs()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, state, ipv4CIDR string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
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
