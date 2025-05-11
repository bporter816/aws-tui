package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
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
		var messages string
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
			utils.DerefString(v.ZoneName, ""),
			utils.DerefString(v.ZoneId, ""),
			utils.TitleCase(string(v.State)),
			messages,
		})
	}
	e.SetData(data)
}
