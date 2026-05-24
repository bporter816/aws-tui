package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewECSTasks(clusterName string, serviceName string, repo *repo.ECS, app *Application) *TableView[model.ECSTask] {
	labels := []string{clusterName, "Tasks"}
	if len(serviceName) > 0 {
		labels = []string{serviceName, "Tasks"}
	}
	return NewTableView(
		"ECS",
		labels,
		[]string{"ID", "GROUP", "LAST STATUS", "DESIRED STATUS", "LAUNCH TYPE", "PLATFORM", "CPU", "MEM"},
		func(ctx context.Context) ([]model.ECSTask, error) {
			return repo.ListTasks(ctx, clusterName, serviceName)
		},
		func(v model.ECSTask) []string {
			var id, group, lastStatus, desiredStatus, platform, cpu, mem string
			if v.TaskArn != nil {
				if a, err := arn.Parse(*v.TaskArn); err == nil {
					id = utils.GetResourceNameFromArn(a)
				} else {
					id = *v.TaskArn
				}
			}
			if v.Group != nil {
				group = *v.Group
			}
			if v.LastStatus != nil {
				lastStatus = utils.AutoCase(*v.LastStatus)
			}
			if v.DesiredStatus != nil {
				desiredStatus = utils.AutoCase(*v.DesiredStatus)
			}
			if v.PlatformFamily != nil {
				platform = *v.PlatformFamily
				if v.PlatformVersion != nil {
					platform += " (" + *v.PlatformVersion + ")"
				}
			}
			if v.Cpu != nil {
				cpu = *v.Cpu + " m"
			}
			if v.Memory != nil {
				mem = *v.Memory + " MB"
			}
			return []string{
				id,
				group,
				lastStatus,
				desiredStatus,
				utils.AutoCase(string(v.LaunchType)),
				platform,
				cpu,
				mem,
			}
		},
		func(v *TableView[model.ECSTask]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TaskArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.TaskArn, app))
					},
				},
			}
		},
	)
}
