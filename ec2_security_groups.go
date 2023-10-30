package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type EC2SecurityGroups struct {
	*ui.Table
	repo *repo.EC2
	app  *Application
}

func NewEC2SecurityGroups(repo *repo.EC2, app *Application) *EC2SecurityGroups {
	e := &EC2SecurityGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"VPC ID",
			"INGRESS RULES",
			"EGRESS RULES",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2SecurityGroups) GetService() string {
	return "EC2"
}

func (e EC2SecurityGroups) GetLabels() []string {
	return []string{"Security Groups"}
}

func (e EC2SecurityGroups) rulesHandler() {
	sgId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupRules(e.repo, sgId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) tagsHandler() {
	sgId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupTags(e.repo, sgId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
			Description: "Rules",
			Action:      e.rulesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2SecurityGroups) Render() {
	model, err := e.repo.ListSecurityGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, description, ingressCount, egressCount string
		vpcId := "-"
		if v.GroupName != nil {
			name = *v.GroupName
		}
		if v.GroupId != nil {
			id = *v.GroupId
		}
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		if v.Description != nil {
			description = *v.Description
		}
		ingressCount = strconv.Itoa(len(v.IpPermissions))
		egressCount = strconv.Itoa(len(v.IpPermissionsEgress))
		data = append(data, []string{
			name,
			id,
			vpcId,
			ingressCount,
			egressCount,
			description,
		})
	}
	e.SetData(data)
}
