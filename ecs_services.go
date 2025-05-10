package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type ECSServices struct {
	*ui.Table
	view.ECS
	repo        *repo.ECS
	app         *Application
	model       []model.ECSService
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

func (e ECSServices) tasksHandler() {
	serviceName, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	tasksView := NewECSTasks(e.clusterName, serviceName, e.repo, e.app)
	e.app.AddAndSwitch(tasksView)
}

func (e ECSServices) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if a := e.model[row-1].ServiceArn; a != nil {
		tagsView := NewTags(e.repo, e.GetService(), *a, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ECSServices) GetKeyActions() []KeyAction {
	return []KeyAction{
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

func (e *ECSServices) Render() {
	model, err := e.repo.ListServices(e.clusterName)
	if err != nil {
		panic(err)
	}
	e.model = model

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
