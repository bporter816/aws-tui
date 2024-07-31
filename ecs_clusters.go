package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"strconv"
)

type ECSClusters struct {
	*ui.Table
	view.ECS
	repo *repo.ECS
	app  *Application
}

func NewECSClusters(repo *repo.ECS, app *Application) *ECSClusters {
	e := &ECSClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"SERVICES",
			"TASKS (P/R)",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ECSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e ECSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ECSClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, status, services, tasks string
		if v.ClusterName != nil {
			name = *v.ClusterName
		}
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
		}
		services = strconv.Itoa(int(v.ActiveServicesCount))
		tasks = strconv.Itoa(int(v.PendingTasksCount)) + "/" + strconv.Itoa(int(v.RunningTasksCount))
		data = append(data, []string{
			name,
			status,
			services,
			tasks,
		})
	}
	e.SetData(data)
}
