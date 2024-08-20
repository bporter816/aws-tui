package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type ELBTargetGroups struct {
	*ui.Table
	view.ELB
	repo  *repo.ELB
	app   *Application
	model []model.ELBTargetGroup
}

func NewELBTargetGroups(repo *repo.ELB, app *Application) *ELBTargetGroups {
	e := &ELBTargetGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"PORT",
			"PROTOCOL",
			"TARGET TYPE",
			"VPC",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ELBTargetGroups) GetLabels() []string {
	return []string{"Target Groups"}
}

func (e ELBTargetGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].TargetGroupArn; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ELBTargetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBTargetGroups) Render() {
	model, err := e.repo.ListTargetGroups()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		// TODO add attached load balancer
		var name, protocol, targetType, vpcId string
		port := "-"
		if v.TargetGroupName != nil {
			name = *v.TargetGroupName
		}
		if v.Port != nil {
			port = strconv.Itoa(int(*v.Port))
		}
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		protocol = string(v.Protocol)
		targetType = utils.AutoCase(string(v.TargetType))
		data = append(data, []string{
			name,
			port,
			protocol,
			targetType,
			vpcId,
		})
	}
	e.SetData(data)
}
