package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewECSClusters(repo *repo.ECS, app *Application) *TableView[model.ECSCluster] {
	return NewTableView(
		"ECS",
		[]string{"Clusters"},
		[]string{"NAME", "STATUS", "SERVICES", "TASKS (P/R)", "CONTAINER INSTANCES"},
		repo.ListClusters,
		func(v model.ECSCluster) []string {
			var status string
			if v.Status != nil {
				status = utils.TitleCase(*v.Status)
			}
			return []string{
				utils.DerefString(v.ClusterName, ""),
				status,
				strconv.Itoa(int(v.ActiveServicesCount)),
				strconv.Itoa(int(v.PendingTasksCount)) + "/" + strconv.Itoa(int(v.RunningTasksCount)),
				strconv.Itoa(int(v.RegisteredContainerInstancesCount)),
			}
		},
		func(v *TableView[model.ECSCluster]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
					Description: "Services",
					Action: func() {
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewECSServices(name, repo, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
					Description: "Tasks",
					Action: func() {
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewECSTasks(name, "", repo, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ClusterArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ClusterArn, app))
					},
				},
			}
		},
	)
}
