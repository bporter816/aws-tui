package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2Instances struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2Instances(repo *repo.EC2, app *Application) *EC2Instances {
	e := &EC2Instances{
		Table: ui.NewTable([]string{
			"NAME",
			"INSTANCE ID",
			"STATE",
			"PUBLIC IP",
			"INSTANCE TYPE",
			"SUBNET ID",
			"KEY NAME",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2Instances) GetLabels() []string {
	return []string{"Instances"}
}

func (e EC2Instances) tagsHandler() {
	instanceId, err := e.GetColSelection("INSTANCE ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), instanceId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2Instances) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
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
		}
		if v.InstanceId != nil {
			id = *v.InstanceId
		}
		if v.State != nil {
			state = utils.TitleCase(string(v.State.Name))
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
