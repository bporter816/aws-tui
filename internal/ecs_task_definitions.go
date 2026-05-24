package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/gdamore/tcell/v2"
)

func NewECSTaskDefinitions(repo *repo.ECS, app *Application) *TableView[string] {
	return NewTableView(
		"ECS",
		[]string{"Task Definitions"},
		[]string{"FAMILY"},
		repo.ListTaskDefinitions,
		func(v string) []string {
			return []string{v}
		},
		func(v *TableView[string]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
					Description: "Revisions",
					Action: func() {
						if family, err := v.GetColSelection("FAMILY"); err == nil {
							app.AddAndSwitch(NewECSTaskDefinitionRevisions(family, repo, app))
						}
					},
				},
			}
		},
	)
}

func init() {
	register("ECS.Task Definitions", func(r Repos, app *Application) Component {
		return NewECSTaskDefinitions(r.ECS, app)
	})
}
