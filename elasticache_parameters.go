package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheParameters struct {
	*ui.Table
	ecClient           *ec.Client
	parameterGroupName string
	app                *Application
}

func NewElasticacheParameters(ecClient *ec.Client, parameterGroupName string, app *Application) *ElasticacheParameters {
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
		ecClient:           ecClient,
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
	pg := ec.NewDescribeCacheParametersPaginator(
		e.ecClient,
		&ec.DescribeCacheParametersInput{
			CacheParameterGroupName: aws.String(e.parameterGroupName),
		},
	)
	var parameters []ecTypes.Parameter
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		parameters = append(parameters, out.Parameters...)
	}

	var data [][]string
	for _, v := range parameters {
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
