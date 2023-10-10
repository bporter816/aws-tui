package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ElasticacheUpdateActions struct {
	*ui.Table
	repo                *repo.Elasticache
	app                 *Application
	cacheClusterIds     []string
	replicationGroupIds []string
	serviceUpdateName   string
}

func NewElasticacheUpdateActions(repo *repo.Elasticache, app *Application, cacheClusterIds []string, replicationGroupIds []string, serviceUpdateName string) *ElasticacheUpdateActions {
	idCol := "UPDATE NAME"
	if serviceUpdateName != "" {
		idCol = "CLUSTER"
	}
	e := &ElasticacheUpdateActions{
		Table: ui.NewTable([]string{
			idCol,
			"STATUS",
			"NODES UPDATED",
		}, 1, 0),
		repo:                repo,
		app:                 app,
		cacheClusterIds:     cacheClusterIds,
		replicationGroupIds: replicationGroupIds,
		serviceUpdateName:   serviceUpdateName,
	}
	return e
}

func (e ElasticacheUpdateActions) GetService() string {
	return "Elasticache"
}

func (e ElasticacheUpdateActions) GetLabels() []string {
	var label string
	if len(e.cacheClusterIds) > 0 {
		label = e.cacheClusterIds[0]
	} else if len(e.replicationGroupIds) > 0 {
		label = e.replicationGroupIds[0]
	} else {
		label = e.serviceUpdateName
	}
	return []string{label, "Update Status"}
}

func (e ElasticacheUpdateActions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheUpdateActions) Render() {
	model, err := e.repo.ListUpdateActions(e.cacheClusterIds, e.replicationGroupIds, e.serviceUpdateName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, status, nodesUpdated string
		if e.serviceUpdateName == "" {
			if v.ServiceUpdateName != nil {
				id = *v.ServiceUpdateName
			}
		} else {
			if v.CacheClusterId != nil {
				id = *v.CacheClusterId
			} else if v.ReplicationGroupId != nil {
				id = *v.ReplicationGroupId
			}
		}
		status = utils.AutoCase(string(v.UpdateActionStatus))
		if v.NodesUpdated != nil {
			nodesUpdated = *v.NodesUpdated
		}
		data = append(data, []string{
			id,
			status,
			nodesUpdated,
		})
	}
	e.SetData(data)
}
