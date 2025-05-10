package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type ECSServices struct {
	*ui.Table
	view.ECS
	repo        *repo.ECS
	app         *Application
	clusterName string
}

func NewECSServices(clusterName string, repo *repo.ECS, app *Application) *ECSServices {
	e := &ECSServices{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"TASKS (P/R)",
			"TASK DEFINITION",
			"TYPE",
		}, 1, 0),
		clusterName: clusterName,
		repo:        repo,
		app:         app,
	}
	return e
}

func (e ECSServices) GetLabels() []string {
	return []string{e.clusterName, "Services"}
}

func (e ECSServices) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ECSServices) Render() {
	model, err := e.repo.ListServices(e.clusterName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
		data = append(data, []string{
			name,
			status,
			tasks,
			taskDefinition,
			scheduling,
		})
	}
	e.SetData(data)
}
