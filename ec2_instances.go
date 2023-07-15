package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type EC2Instances struct {
	*ui.Table
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2Instances(ec2Client *ec2.Client, app *Application) *EC2Instances {
	e := &EC2Instances{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"STATE",
			"INSTANCE TYPE",
			"SUBNET ID",
			"KEY NAME",
		}, 1, 0),
		ec2Client: ec2Client,
		app:       app,
	}
	return e
}

func (e EC2Instances) GetService() string {
	return "EC2"
}

func (e EC2Instances) GetLabels() []string {
	return []string{"Instances"}
}

func (e EC2Instances) tagsHandler() {
	instanceId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2InstanceTags(e.ec2Client, instanceId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2Instances) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2Instances) Render() {
	pg := ec2.NewDescribeInstancesPaginator(
		e.ec2Client,
		&ec2.DescribeInstancesInput{},
	)
	var instances []ec2Types.Instance
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		for _, v := range out.Reservations {
			instances = append(instances, v.Instances...)
		}
	}

	var data [][]string
	for _, v := range instances {
		var name, id, state, instanceType, subnetId, keyName string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
		}
		if v.InstanceId != nil {
			id = *v.InstanceId
		}
		if v.State != nil {
			state = string(v.State.Name)
		}
		instanceType = string(v.InstanceType)
		if v.SubnetId != nil {
			subnetId = *v.SubnetId
		}
		if v.KeyName != nil {
			keyName = *v.KeyName
		}
		data = append(data, []string{
			name,
			id,
			state,
			instanceType,
			subnetId,
			keyName,
		})
	}
	e.SetData(data)
}
