package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ElastiCacheParameters struct {
	*ui.Table
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

func (e ElastiCacheParameters) GetService() string {
	return "ElastiCache"
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
		var name, allowedValues, value, dataType, isModifiable, source, description string
		if v.ParameterName != nil {
			name = *v.ParameterName
		}
		if v.AllowedValues != nil {
			allowedValues = *v.AllowedValues
		}
		if v.IsModifiable != nil {
			isModifiable = utils.BoolToString(*v.IsModifiable, "Yes", "No")
		}
		if v.ParameterValue != nil {
			value = *v.ParameterValue
		}
		if v.DataType != nil {
			dataType = *v.DataType
		}
		var changeType = string(v.ChangeType)
		if v.Source != nil {
			source = *v.Source
		}
		if v.Description != nil {
			description = *v.Description
		}
		data = append(data, []string{
			name,
			allowedValues,
			value,
			dataType,
			isModifiable,
			changeType,
			source,
			description,
		})
	}
	e.SetData(data)
}
