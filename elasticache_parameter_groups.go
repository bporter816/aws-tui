package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheParameterGroups struct {
	*ui.Table
	view.ElastiCache
	repo *repo.ElastiCache
	app  *Application
}

func NewElastiCacheParameterGroups(repo *repo.ElastiCache, app *Application) *ElastiCacheParameterGroups {
	e := &ElastiCacheParameterGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"FAMILY",
			"DESCRIPTION",
			"GLOBAL",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElastiCacheParameterGroups) GetLabels() []string {
	return []string{"Parameter Groups"}
}

func (e ElastiCacheParameterGroups) viewParametersHandler() {
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	parametersView := NewElastiCacheParameters(e.repo, name, e.app)
	e.app.AddAndSwitch(parametersView)
}

func (e ElastiCacheParameterGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "View Parameters",
			Action:      e.viewParametersHandler,
		},
	}
}

func (e ElastiCacheParameterGroups) Render() {
	model, err := e.repo.ListParameterGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var isGlobal string
		if v.IsGlobal != nil {
			isGlobal = utils.BoolToString(*v.IsGlobal, "Yes", "No")
		}
		data = append(data, []string{
			utils.DerefString(v.CacheParameterGroupName, ""),
			utils.DerefString(v.CacheParameterGroupFamily, ""),
			utils.DerefString(v.Description, ""),
			isGlobal,
		})
	}
	e.SetData(data)
}
