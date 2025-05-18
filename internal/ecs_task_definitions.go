package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ECSTaskDefinitions struct {
	*ui.Table
	view.ECS
	repo  *repo.ECS
	app   *Application
	model []string
}

func NewECSTaskDefinitions(repo *repo.ECS, app *Application) *ECSTaskDefinitions {
	e := &ECSTaskDefinitions{
		Table: ui.NewTable([]string{
			"FAMILY",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ECSTaskDefinitions) GetLabels() []string {
	return []string{"Task Definitions"}
}

func (e ECSTaskDefinitions) revisionsHandler() {
	if family, err := e.GetColSelection("FAMILY"); err == nil {
		revisionsView := NewECSTaskDefinitionRevisions(family, e.repo, e.app)
		e.app.AddAndSwitch(revisionsView)
	}
}

func (e ECSTaskDefinitions) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
			Description: "Revisions",
			Action:      e.revisionsHandler,
		},
	}
}

func (e *ECSTaskDefinitions) Render() {
	model, err := e.repo.ListTaskDefinitions()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{v})
	}
	e.SetData(data)
}
