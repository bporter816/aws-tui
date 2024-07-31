package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ECSClusters struct {
	*ui.Table
	view.ECS
	repo  *repo.ECS
	app   *Application
	model []model.ECSCluster
}

func NewECSClusters(repo *repo.ECS, app *Application) *ECSClusters {
	e := &ECSClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"SERVICES",
			"TASKS (P/R)",
			"CONTAINER INSTANCES",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ECSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e ECSClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ClusterArn; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ECSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ECSClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var name, status, services, tasks, containerInstances string
		if v.ClusterName != nil {
			name = *v.ClusterName
		}
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
		}
		services = strconv.Itoa(int(v.ActiveServicesCount))
		tasks = strconv.Itoa(int(v.PendingTasksCount)) + "/" + strconv.Itoa(int(v.RunningTasksCount))
		containerInstances = strconv.Itoa(int(v.RegisteredContainerInstancesCount))
		data = append(data, []string{
			name,
			status,
			services,
			tasks,
			containerInstances,
		})
	}
	e.SetData(data)
}
