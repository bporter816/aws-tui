package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/ui"
)

type EC2VPCTags struct {
	*ui.Table
	ec2Client *ec2.Client
	vpcId     string
	app       *Application
}

func NewEC2VPCTags(ec2Client *ec2.Client, vpcId string, app *Application) *EC2VPCTags {
	e := &EC2VPCTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client: ec2Client,
		vpcId:     vpcId,
		app:       app,
	}
	return e
}

func (e EC2VPCTags) GetService() string {
	return "EC2"
}

func (e EC2VPCTags) GetLabels() []string {
	return []string{e.vpcId, "Tags"}
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
