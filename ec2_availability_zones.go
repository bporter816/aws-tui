package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type EC2AvailabilityZones struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2AvailabilityZones(repo *repo.EC2, app *Application) *EC2AvailabilityZones {
	e := &EC2AvailabilityZones{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"STATE",
			"MESSAGES",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2AvailabilityZones) GetLabels() []string {
	return []string{"Availability Zones"}
}

func (e EC2AvailabilityZones) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2AvailabilityZones) Render() {
	model, err := e.repo.ListAvailabilityZones()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, state, messages string
		if v.ZoneName != nil {
			name = *v.ZoneName
		}
		if v.ZoneId != nil {
			id = *v.ZoneId
		}
		state = utils.TitleCase(string(v.State))
		if len(v.Messages) > 0 {
			for _, m := range v.Messages {
				if m.Message != nil {
					messages += ", " + *m.Message
				}
			}
			messages = messages[2:]
		} else {
			messages = "-"
		}
		data = append(data, []string{
			name,
			id,
			state,
			messages,
		})
	}
	e.SetData(data)
}
