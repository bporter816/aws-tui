package internal

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewECSServices(clusterName string, repo *repo.ECS, app *Application) *TableView[model.ECSService] {
	return NewTableView(
		"ECS",
		[]string{clusterName, "Services"},
		[]string{"NAME", "STATUS", "TASKS (P/R)", "TASK DEFINITION", "TYPE"},
		func() ([]model.ECSService, error) {
			return repo.ListServices(clusterName)
		},
		func(v model.ECSService) []string {
			var name, status, tasks, taskDefinition, scheduling string
			if v.ServiceName != nil {
				name = *v.ServiceName
			}
			if v.Status != nil {
				status = utils.AutoCase(*v.Status)
			}
			tasks = strconv.Itoa(int(v.PendingCount)) + "/" + strconv.Itoa(int(v.RunningCount))
			if v.TaskDefinition != nil {
				if a, err := arn.Parse(*v.TaskDefinition); err == nil {
					taskDefinition = utils.GetResourceNameFromArn(a)
				}
			}
			scheduling = utils.AutoCase(string(v.SchedulingStrategy))
			return []string{
				name,
				status,
				tasks,
				taskDefinition,
				scheduling,
			}
		},
		func(v *TableView[model.ECSService]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
					Description: "Tasks",
					Action: func() {
						serviceName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewECSTasks(clusterName, serviceName, repo, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ServiceArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ServiceArn, app))
					},
				},
			}
		},
	)
}
