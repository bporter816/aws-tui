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

func (e ECSClusters) servicesHandler() {
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	servicesView := NewECSServices(name, e.repo, e.app)
	e.app.AddAndSwitch(servicesView)
}

func (e ECSClusters) tasksHandler() {
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	tasksView := NewECSTasks(name, e.repo, e.app)
	e.app.AddAndSwitch(tasksView)
}

func (e ECSClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if a := e.model[row-1].ClusterArn; a != nil {
		tagsView := NewTags(e.repo, e.GetService(), *a, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ECSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Services",
			Action:      e.servicesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tasks",
			Action:      e.tasksHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var status string
		if v.Status != nil {
			status = utils.TitleCase(*v.Status)
		}
		data = append(data, []string{
			utils.DerefString(v.ClusterName, ""),
			status,
			strconv.Itoa(int(v.ActiveServicesCount)),
			strconv.Itoa(int(v.PendingTasksCount)) + "/" + strconv.Itoa(int(v.RunningTasksCount)),
			strconv.Itoa(int(v.RegisteredContainerInstancesCount)),
		})
	}
	e.SetData(data)
}
