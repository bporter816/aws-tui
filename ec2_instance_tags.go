package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/ui"
)

type EC2InstanceTags struct {
	*ui.Table
	ec2Client  *ec2.Client
	instanceId string
	app        *Application
}

func NewEC2InstanceTags(ec2Client *ec2.Client, instanceId string, app *Application) *EC2InstanceTags {
	e := &EC2InstanceTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ec2Client:  ec2Client,
		instanceId: instanceId,
		app:        app,
	}
	return e
}

func (e EC2InstanceTags) GetService() string {
	return "EC2"
}

func (e EC2InstanceTags) GetLabels() []string {
	return []string{e.instanceId, "Tags"}
}

func (e EC2InstanceTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2InstanceTags) Render() {
	out, err := e.ec2Client.DescribeInstances(
		context.TODO(),
		&ec2.DescribeInstancesInput{
			InstanceIds: []string{e.instanceId},
		},
	)
	if err != nil {
		panic(err)
	}

	if len(out.Reservations) != 1 {
		panic("should get exactly 1 reservation")
	}
	if len(out.Reservations[0].Instances) != 1 {
		panic("should get exactly 1 instance")
	}
	var data [][]string
	for _, v := range out.Reservations[0].Instances[0].Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
