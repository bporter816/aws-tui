package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheParameters struct {
	*ui.Table
	repo               *repo.Elasticache
	parameterGroupName string
	app                *Application
}

func NewElasticacheParameters(repo *repo.Elasticache, parameterGroupName string, app *Application) *ElasticacheParameters {
	e := &ElasticacheParameters{
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

func (e ElasticacheParameters) GetService() string {
	return "Elasticache"
}

func (e ElasticacheParameters) GetLabels() []string {
	return []string{e.parameterGroupName, "Parameters"}
}

func (e ElasticacheParameters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheParameters) Render() {
	model, err := e.repo.ListParameters(e.parameterGroupName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		name, allowedValues, value, dataType, isModifiable, source, description := "", "", "", "", "No", "", ""
		if v.ParameterName != nil {
			name = *v.ParameterName
		}
		if v.AllowedValues != nil {
			allowedValues = *v.AllowedValues
		}
		if v.IsModifiable {
			isModifiable = "Yes"
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
