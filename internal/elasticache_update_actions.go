package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type ElastiCacheUpdateActions struct {
	*ui.Table
	view.ElastiCache
	repo                *repo.ElastiCache
	app                 *Application
	cacheClusterIds     []string
	replicationGroupIds []string
	serviceUpdateName   string
}

func NewElastiCacheUpdateActions(repo *repo.ElastiCache, app *Application, cacheClusterIds []string, replicationGroupIds []string, serviceUpdateName string) *ElastiCacheUpdateActions {
	idCol := "UPDATE NAME"
	if serviceUpdateName != "" {
		idCol = "CLUSTER"
	}
	e := &ElastiCacheUpdateActions{
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

func (e ElastiCacheUpdateActions) GetLabels() []string {
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

func (e ElastiCacheUpdateActions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElastiCacheUpdateActions) Render() {
	model, err := e.repo.ListUpdateActions(e.cacheClusterIds, e.replicationGroupIds, e.serviceUpdateName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id string
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
		data = append(data, []string{
			id,
			utils.AutoCase(string(v.UpdateActionStatus)),
			utils.DerefString(v.NodesUpdated, ""),
		})
	}
	e.SetData(data)
}
