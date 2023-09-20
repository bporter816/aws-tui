package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type ElasticacheParameterGroups struct {
	*ui.Table
	repo     *repo.Elasticache
	ecClient *ec.Client
	app      *Application
}

func NewElasticacheParameterGroups(repo *repo.Elasticache, ecClient *ec.Client, app *Application) *ElasticacheParameterGroups {
	e := &ElasticacheParameterGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"FAMILY",
			"DESCRIPTION",
			"GLOBAL",
		}, 1, 0),
		repo:     repo,
		ecClient: ecClient,
		app:      app,
	}
	return e
}

func (e ElasticacheParameterGroups) GetService() string {
	return "Elasticache"
}

func (e ElasticacheParameterGroups) GetLabels() []string {
	return []string{"Parameter Groups"}
}

func (e ElasticacheParameterGroups) viewParametersHandler() {
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	parametersView := NewElasticacheParameters(e.repo, name, e.app)
	e.app.AddAndSwitch(parametersView)
}

func (e ElasticacheParameterGroups) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "View Parameters",
			Action:      e.viewParametersHandler,
		},
	}
}

func (e ElasticacheParameterGroups) Render() {
	pg := ec.NewDescribeCacheParameterGroupsPaginator(
		e.ecClient,
		&ec.DescribeCacheParameterGroupsInput{},
	)
	var parameterGroups []ecTypes.CacheParameterGroup
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		parameterGroups = append(parameterGroups, out.CacheParameterGroups...)
	}

	var data [][]string
	for _, v := range parameterGroups {
		name, family, description, isGlobal := "", "", "", "No"
		if v.CacheParameterGroupName != nil {
			name = *v.CacheParameterGroupName
		}
		if v.CacheParameterGroupFamily != nil {
			family = *v.CacheParameterGroupFamily
		}
		if v.Description != nil {
			description = *v.Description
		}
		if v.IsGlobal {
			isGlobal = "Yes"
		}
		data = append(data, []string{
			name,
			family,
			description,
			isGlobal,
		})
	}
	e.SetData(data)
}
