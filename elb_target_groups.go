package main

import (
	"context"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"strconv"
)

type ELBTargetGroups struct {
	*Table
	elbClient *elb.Client
	app       *Application
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

func (e ELBTargetGroups) GetName() string {
	return "ELB | Target Groups"
}

func (e ELBTargetGroups) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTargetGroups) Render() {
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
	for _, v := range targetGroups {
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
