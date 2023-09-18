package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type EC2Instances struct {
	*ui.Table
	repo      *repo.EC2
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2Instances(repo *repo.EC2, ec2Client *ec2.Client, app *Application) *EC2Instances {
	e := &EC2Instances{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"STATE",
			"PUBLIC IP",
			"INSTANCE TYPE",
			"SUBNET ID",
			"KEY NAME",
		}, 1, 0),
		repo:      repo,
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
	model, err := e.repo.ListInstances()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, state, publicIP, instanceType, subnetId, keyName string
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
		if v.PublicIpAddress != nil {
			publicIP = *v.PublicIpAddress
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
			publicIP,
			instanceType,
			subnetId,
			keyName,
		})
	}
	e.SetData(data)
}
