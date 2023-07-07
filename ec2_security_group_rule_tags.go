package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2SecurityGroupRuleTags struct {
	*Table
	ec2Client *ec2.Client
	sgId      string
	ruleId    string
	app       *Application
}

func NewEC2SecurityGroupRuleTags(ec2Client *ec2.Client, sgId string, ruleId string, app *Application) *EC2SecurityGroupRuleTags {
	e := &EC2SecurityGroupRuleTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		sgId:      sgId,
		ruleId:    ruleId,
		app:       app,
	}
	return e
}

func (e EC2SecurityGroupRuleTags) GetName() string {
	return fmt.Sprintf("EC2 | Security Groups | %v | Rules | %v | Tags", e.sgId, e.ruleId)
}

func (e EC2SecurityGroupRuleTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2SecurityGroupRuleTags) Render() {
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
		panic("should get exactly 1 rule")
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
