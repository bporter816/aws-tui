package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/ui"
)

type EC2SecurityGroupTags struct {
	*ui.Table
	ec2Client *ec2.Client
	sgId      string
	app       *Application
}

func NewEC2SecurityGroupTags(ec2Client *ec2.Client, sgId string, app *Application) *EC2SecurityGroupTags {
	e := &EC2SecurityGroupTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		sgId:      sgId,
		app:       app,
	}
	return e
}

func (e EC2SecurityGroupTags) GetService() string {
	return "EC2"
}

func (e EC2SecurityGroupTags) GetLabels() []string {
	return []string{e.sgId, "Tags"}
}

func (e EC2SecurityGroupTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2SecurityGroupTags) Render() {
	out, err := e.ec2Client.DescribeSecurityGroups(
		context.TODO(),
		&ec2.DescribeSecurityGroupsInput{
			GroupIds: []string{e.sgId},
		},
	)
	if err != nil {
		panic(err)
	}

	if len(out.SecurityGroups) != 1 {
		panic("should get exactly 1 security group")
	}
	var data [][]string
	for _, v := range out.SecurityGroups[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
