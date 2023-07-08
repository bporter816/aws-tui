package main

import (
	"context"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ELBTargetGroups struct {
	*Table
	elbClient *elb.Client
	app       *Application
	arns      []string
}

func NewELBTargetGroups(elbClient *elb.Client, app *Application) *ELBTargetGroups {
	e := &ELBTargetGroups{
		Table: NewTable([]string{
			"NAME",
			"PORT",
			"PROTOCOL",
			"TARGET TYPE",
			"VPC",
		}, 1, 0),
		elbClient: elbClient,
		app:       app,
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
	tagsView := NewELBTags(e.elbClient, ELBResourceTypeTargetGroup, e.arns[row-1], name, e.app)
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
	pg := elb.NewDescribeTargetGroupsPaginator(
		e.elbClient,
		&elb.DescribeTargetGroupsInput{},
	)
	var targetGroups []elbTypes.TargetGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		targetGroups = append(targetGroups, out.TargetGroups...)
	}

	var data [][]string
	e.arns = make([]string, len(targetGroups))
	for i, v := range targetGroups {
		e.arns[i] = *v.TargetGroupArn
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
