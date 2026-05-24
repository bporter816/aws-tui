package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewSQSQueues(repo *repo.SQS, app *Application) *TableView[model.SQSQueue] {
	return NewTableView(
		"SQS",
		[]string{"Queues"},
		[]string{"NAME", "TYPE"},
		repo.ListQueues,
		func(v model.SQSQueue) []string {
			return []string{
				v.Name,
				utils.BoolToString(v.IsFifo, "FIFO", "Standard"),
			}
		},
		func(v *TableView[model.SQSQueue]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Access Policy",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSQSAccessPolicy(repo, item.QueueUrl, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), item.QueueUrl, app))
					},
				},
			}
		},
	)
}
