package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type ElastiCacheTags struct {
	*ui.Table
	repo         *repo.ElastiCache
	resourceArn  string
	resourceName string
	app          *Application
}

type ElastiCacheResourceType string

func NewElastiCacheTags(repo *repo.ElastiCache, resourceArn string, resourceName string, app *Application) *ElastiCacheTags {
	e := &ElastiCacheTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:         repo,
		resourceArn:  resourceArn,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ElastiCacheTags) GetService() string {
	return "ElastiCache"
}

func (e ElastiCacheTags) GetLabels() []string {
	return []string{e.resourceName, "Tags"}
}

func (e ElastiCacheTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElastiCacheTags) Render() {
	model, err := e.repo.ListTags(e.resourceArn)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	e.SetData(data)
}
