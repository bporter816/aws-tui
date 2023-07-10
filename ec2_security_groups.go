package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type EC2SecurityGroups struct {
	*Table
	ec2Client *ec2.Client
	app       *Application
	filters   []ec2Types.Filter
}

func NewEC2SecurityGroups(ec2Client *ec2.Client, app *Application, filters ...ec2Types.Filter) *EC2SecurityGroups {
	e := &EC2SecurityGroups{
		Table: NewTable([]string{
			"NAME",
			"ID",
			"VPC ID",
			"INGRESS RULES",
			"EGRESS RULES",
			"DESCRIPTION",
		}, 1, 0),
		ec2Client: ec2Client,
		app:       app,
		filters:   filters,
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
	tagsView := NewEC2SecurityGroupRules(e.ec2Client, sgId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) tagsHandler() {
	sgId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	rulesView := NewEC2SecurityGroupTags(e.ec2Client, sgId, e.app)
	e.app.AddAndSwitch(rulesView)
}

func (e EC2SecurityGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
			Description: "Rules",
			Action:      e.rulesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2SecurityGroups) Render() {
	pg := ec2.NewDescribeSecurityGroupsPaginator(
		e.ec2Client,
		&ec2.DescribeSecurityGroupsInput{
			Filters: e.filters,
		},
	)
	var securityGroups []ec2Types.SecurityGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		securityGroups = append(securityGroups, out.SecurityGroups...)
	}

	var data [][]string
	for _, v := range securityGroups {
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
