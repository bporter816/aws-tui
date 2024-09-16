package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type ElastiCacheParameters struct {
	*ui.Table
	view.ElastiCache
	repo               *repo.ElastiCache
	parameterGroupName string
	app                *Application
}

func NewElastiCacheParameters(repo *repo.ElastiCache, parameterGroupName string, app *Application) *ElastiCacheParameters {
	e := &ElastiCacheParameters{
		Table: ui.NewTable([]string{
			"NAME",
			"ALLOWED VALUES",
			"VALUE",
			"TYPE",
			"MODIFIABLE",
			"CHANGE TYPE",
			"SOURCE",
			"DESCRIPTION",
		}, 1, 0),
		repo:               repo,
		parameterGroupName: parameterGroupName,
		app:                app,
	}
	return e
}

func (e ElastiCacheParameters) GetLabels() []string {
	return []string{e.parameterGroupName, "Parameters"}
}

func (e ElastiCacheParameters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElastiCacheParameters) Render() {
	model, err := e.repo.ListParameters(e.parameterGroupName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var isModifiable string
		if v.IsModifiable != nil {
			isModifiable = utils.BoolToString(*v.IsModifiable, "Yes", "No")
		}
		data = append(data, []string{
			utils.DerefString(v.ParameterName, ""),
			utils.DerefString(v.AllowedValues, ""),
			utils.DerefString(v.ParameterValue, ""),
			utils.DerefString(v.DataType, ""),
			isModifiable,
			string(v.ChangeType),
			utils.DerefString(v.Source, ""),
			utils.DerefString(v.Description, ""),
		})
	}
	e.SetData(data)
}
