package main

import (
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ELBTargetGroups struct {
	*ui.Table
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

func (e ELBTargetGroups) GetService() string {
	return "ELB"
}

func (e ELBTargetGroups) GetLabels() []string {
	return []string{"Target Groups"}
}

func (e ELBTargetGroups) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if e.model[row-1].TargetGroupArn == nil {
		return
	}
	tagsView := NewELBTags(e.repo, ELBResourceTypeTargetGroup, *e.model[row-1].TargetGroupArn, name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ELBTargetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		targetType = formatTargetType(v.TargetType)
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

func formatTargetType(e elbTypes.TargetTypeEnum) string {
	switch e {
	case elbTypes.TargetTypeEnumInstance:
		return "Instance"
	case elbTypes.TargetTypeEnumIp:
		return "IP"
	case elbTypes.TargetTypeEnumLambda:
		return "Lambda"
	case elbTypes.TargetTypeEnumAlb:
		return "ALB"
	default:
		return "<unknown>"
	}
}
