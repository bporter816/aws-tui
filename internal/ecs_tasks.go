package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ECSTasks struct {
	*ui.Table
	view.ECS
	repo        *repo.ECS
	app         *Application
	model       []model.ECSTask
	clusterName string
	serviceName string
}

func NewECSTasks(clusterName string, serviceName string, repo *repo.ECS, app *Application) *ECSTasks {
	e := &ECSTasks{
		Table: ui.NewTable([]string{
			"ID",
			"GROUP",
			"LAST STATUS",
			"DESIRED STATUS",
			"LAUNCH TYPE",
			"PLATFORM",
			"CPU",
			"MEM",
		}, 1, 0),
		clusterName: clusterName,
		serviceName: serviceName,
		repo:        repo,
		app:         app,
	}
	return e
}

func (e ECSTasks) GetLabels() []string {
	if len(e.serviceName) == 0 {
		return []string{e.clusterName, "Tasks"}
	}
	return []string{e.serviceName, "Tasks"}
}

func (e ECSTasks) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if a := e.model[row-1].TaskArn; a != nil {
		tagsView := NewTags(e.repo, e.GetService(), *a, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ECSTasks) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ECSTasks) Render() {
	model, err := e.repo.ListTasks(e.clusterName, e.serviceName)
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		var id, group, lastStatus, desiredStatus, platform, cpu, mem string
		if v.TaskArn != nil {
			a, err := arn.Parse(*v.TaskArn)
			if err != nil {
				panic(err)
			}
			id = utils.GetResourceNameFromArn(a)
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
		data = append(data, []string{
			id,
			group,
			lastStatus,
			desiredStatus,
			utils.AutoCase(string(v.LaunchType)),
			platform,
			cpu,
			mem,
		})
	}
	e.SetData(data)
}
