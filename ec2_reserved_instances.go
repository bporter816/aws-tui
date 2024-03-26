package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2ReservedInstances struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2ReservedInstances(repo *repo.EC2, app *Application) *EC2ReservedInstances {
	e := &EC2ReservedInstances{
		Table: ui.NewTable([]string{
			"ID",
			"INSTANCE TYPE",
			"SCOPE",
			"STATE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2ReservedInstances) GetLabels() []string {
	return []string{"Reserved Instances"}
}

func (e EC2ReservedInstances) tagsHandler() {
	id, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), id, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2ReservedInstances) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2ReservedInstances) Render() {
	model, err := e.repo.ListReservedInstances(nil)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id string
		if v.ReservedInstancesId != nil {
			id = *v.ReservedInstancesId
		}
		data = append(data, []string{
			id,
			string(v.InstanceType),
			string(v.Scope),
			utils.AutoCase(string(v.State)),
		})
	}
	e.SetData(data)
}
