package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ECSTaskDefinitionRevisions struct {
	*ui.Table
	view.ECS
	repo   *repo.ECS
	app    *Application
	family string
	model  []string
}

func NewECSTaskDefinitionRevisions(family string, repo *repo.ECS, app *Application) *ECSTaskDefinitionRevisions {
	e := &ECSTaskDefinitionRevisions{
		Table: ui.NewTable([]string{
			"REVISION",
		}, 1, 0),
		family: family,
		repo:   repo,
		app:    app,
	}
	return e
}

func (e ECSTaskDefinitionRevisions) GetLabels() []string {
	return []string{e.family, "Revisions"}
}

func (e ECSTaskDefinitionRevisions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e *ECSTaskDefinitionRevisions) Render() {
	model, err := e.repo.ListTaskDefinitionRevisions(e.family)
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		a, err := arn.Parse(v)
		if err != nil {
			panic(err)
		}
		data = append(data, []string{
			utils.GetResourceNameFromArn(a),
		})
	}
	e.SetData(data)
}
