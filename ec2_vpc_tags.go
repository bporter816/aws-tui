package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2VPCTags struct {
	*Table
	ec2Client *ec2.Client
	vpcId     string
	app       *Application
}

func NewEC2VPCTags(ec2Client *ec2.Client, vpcId string, app *Application) *EC2VPCTags {
	e := &EC2VPCTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		vpcId:     vpcId,
		app:       app,
	}
	return e
}

func (e EC2VPCTags) GetName() string {
	return fmt.Sprintf("EC2 | VPCs | %v | Tags", e.vpcId)
}

func (e EC2VPCTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2VPCTags) Render() {
	out, err := e.ec2Client.DescribeVpcs(
		context.TODO(),
		&ec2.DescribeVpcsInput{
			VpcIds: []string{e.vpcId},
		},
	)
	if err != nil {
		panic(err)
	}

	if len(out.Vpcs) != 1 {
		panic("should get exactly 1 vpc")
	}
	var data [][]string
	for _, v := range out.Vpcs[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
