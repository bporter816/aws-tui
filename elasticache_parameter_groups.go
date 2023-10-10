package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheParameterGroups struct {
	*ui.Table
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

func (e ElastiCacheParameterGroups) GetService() string {
	return "ElastiCache"
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
		KeyAction{
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
		var name, family, description, isGlobal string
		if v.CacheParameterGroupName != nil {
			name = *v.CacheParameterGroupName
		}
		if v.CacheParameterGroupFamily != nil {
			family = *v.CacheParameterGroupFamily
		}
		if v.Description != nil {
			description = *v.Description
		}
		isGlobal = utils.BoolToString(v.IsGlobal, "Yes", "No")
		data = append(data, []string{
			name,
			family,
			description,
			isGlobal,
		})
	}
	e.SetData(data)
}
