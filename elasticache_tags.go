package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheTags struct {
	*ui.Table
	repo         *repo.Elasticache
	resourceType ElasticacheResourceType
	resourceArn  string
	resourceName string
	app          *Application
}

type ElasticacheResourceType string

const (
	ElasticacheResourceTypeCluster      ElasticacheResourceType = "Clusters"
	ElasticacheResourceTypeReservedNode ElasticacheResourceType = "Reserved Nodes"
	ElasticacheResourceTypeSnapshot     ElasticacheResourceType = "Snapshots"
)

func NewElasticacheTags(repo *repo.Elasticache, resourceType ElasticacheResourceType, resourceArn string, resourceName string, app *Application) *ElasticacheTags {
	e := &ElasticacheTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:         repo,
		resourceType: resourceType,
		resourceArn:  resourceArn,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ElasticacheTags) GetService() string {
	return "Elasticache"
}

func (e ElasticacheTags) GetLabels() []string {
	return []string{e.resourceName, "Tags"}
}

func (e ElasticacheTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheTags) Render() {
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
