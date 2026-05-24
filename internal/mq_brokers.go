package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewMQBrokers(repo *repo.MQ, app *Application) *TableView[model.MQBroker] {
	return NewTableView(
		"MQ",
		[]string{"Brokers"},
		[]string{"NAME", "STATE", "ENGINE", "MODE", "INSTANCE TYPE", "CREATED"},
		repo.ListBrokers,
		func(v model.MQBroker) []string {
			var created string
			if v.Created != nil {
				created = v.Created.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.BrokerName, ""),
				utils.AutoCase(string(v.BrokerState)),
				utils.AutoCase(string(v.EngineType)),
				string(v.DeploymentMode), // TODO format this
				utils.DerefString(v.HostInstanceType, ""),
				created,
			}
		},
		func(v *TableView[model.MQBroker]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.BrokerArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.BrokerArn, app))
					},
				},
			}
		},
	)
}
