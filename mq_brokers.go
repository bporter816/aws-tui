package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type MQBrokers struct {
	*ui.Table
	view.MQ
	repo  *repo.MQ
	app   *Application
	model []model.MQBroker
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

func (m MQBrokers) tagsHandler() {
	row, err := m.GetRowSelection()
	if err != nil {
		return
	}
	if arn := m.model[row-1].BrokerArn; arn != nil {
		tagsView := NewTags(m.repo, m.GetService(), *arn, m.app)
		m.app.AddAndSwitch(tagsView)
	}
}

func (m MQBrokers) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      m.tagsHandler,
		},
	}
}

func (m *MQBrokers) Render() {
	model, err := m.repo.ListBrokers()
	if err != nil {
		panic(err)
	}
	m.model = model

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
