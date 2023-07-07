package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gdamore/tcell/v2"
)

type EC2SecurityGroups struct {
	*Table
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2SecurityGroups(ec2Client *ec2.Client, app *Application) *EC2SecurityGroups {
	e := &EC2SecurityGroups{
		Table: NewTable([]string{
			"NAME",
			"ID",
			"VPC ID",
			"DESCRIPTION",
		}, 1, 0),
		ec2Client: ec2Client,
		app:       app,
	}
	return e
}

func (e EC2SecurityGroups) GetName() string {
	return "EC2 | Security Groups"
}

func (e EC2SecurityGroups) tagsHandler() {
	vpcId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2SecurityGroupTags(e.ec2Client, vpcId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2SecurityGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
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
		&ec2.DescribeSecurityGroupsInput{},
	)
	var vpcs []ec2Types.SecurityGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		vpcs = append(vpcs, out.SecurityGroups...)
	}

	var data [][]string
	for _, v := range vpcs {
		var name, id, description string
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
		data = append(data, []string{
			name,
			id,
			vpcId,
			description,
		})
	}
	e.SetData(data)
}
