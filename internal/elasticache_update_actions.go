package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewElastiCacheUpdateActions(repo *repo.ElastiCache, app *Application, cacheClusterIds []string, replicationGroupIds []string, serviceUpdateName string) *TableView[model.ElastiCacheUpdateAction] {
	idCol := "UPDATE NAME"
	if serviceUpdateName != "" {
		idCol = "CLUSTER"
	}
	var label string
	if len(cacheClusterIds) > 0 {
		label = cacheClusterIds[0]
	} else if len(replicationGroupIds) > 0 {
		label = replicationGroupIds[0]
	} else {
		label = serviceUpdateName
	}
	return NewTableView(
		"ElastiCache",
		[]string{label, "Update Status"},
		[]string{idCol, "STATUS", "NODES UPDATED"},
		func() ([]model.ElastiCacheUpdateAction, error) {
			return repo.ListUpdateActions(cacheClusterIds, replicationGroupIds, serviceUpdateName)
		},
		func(v model.ElastiCacheUpdateAction) []string {
			var id string
			if serviceUpdateName == "" {
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
			return []string{
				id,
				utils.AutoCase(string(v.UpdateActionStatus)),
				utils.DerefString(v.NodesUpdated, ""),
			}
		},
		nil,
	)
}
