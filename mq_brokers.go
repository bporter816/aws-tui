package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type MQBrokers struct {
	*ui.Table
	view.MQ
	repo *repo.MQ
	app  *Application
}

func NewMQBrokers(repo *repo.MQ, app *Application) *MQBrokers {
	m := &MQBrokers{
		Table: ui.NewTable([]string{
			"NAME",
			"STATE",
			"ENGINE",
			"MODE",
			"INSTANCE TYPE",
			"CREATED",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return m
}

func (m MQBrokers) GetLabels() []string {
	return []string{"Brokers"}
}

func (m MQBrokers) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (m MQBrokers) Render() {
	model, err := m.repo.ListBrokers()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var created string
		if v.Created != nil {
			created = v.Created.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.BrokerName, ""),
			utils.AutoCase(string(v.BrokerState)),
			utils.AutoCase(string(v.EngineType)),
			string(v.DeploymentMode), // TODO format this
			utils.DerefString(v.HostInstanceType, ""),
			created,
		})
	}
	m.SetData(data)
}
