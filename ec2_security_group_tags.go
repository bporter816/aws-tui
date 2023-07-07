package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2SecurityGroupTags struct {
	*Table
	ec2Client *ec2.Client
	sgId      string
	app       *Application
}

func NewEC2SecurityGroupTags(ec2Client *ec2.Client, sgId string, app *Application) *EC2SecurityGroupTags {
	e := &EC2SecurityGroupTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		sgId:      sgId,
		app:       app,
	}
	return e
}

func (e EC2SecurityGroupTags) GetName() string {
	return fmt.Sprintf("EC2 | Security Groups | %v | Tags", e.sgId)
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
		panic("should get exactly 1 vpc")
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
