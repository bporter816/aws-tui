package internal

import (
	"strconv"

	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheClusters(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheCluster] {
	return NewTableView(
		"ElastiCache",
		[]string{"Clusters"},
		[]string{"ID", "STATUS", "ENGINE", "VERSION", "NODE TYPE", "CLUSTER MODE", "SHARDS", "NODES"},
		repo.ListClusters,
		func(v model.ElastiCacheCluster) []string {
			if vv := v.CacheCluster; vv != nil {
				var clusterMode string = "-"
				if *vv.Engine == "redis" {
					clusterMode = string(ecTypes.ClusterModeDisabled)
				}
				return []string{
					*vv.CacheClusterId,
					utils.TitleCase(*vv.CacheClusterStatus),
					utils.TitleCase(*vv.Engine),
					*vv.EngineVersion,
					*vv.CacheNodeType,
					utils.TitleCase(clusterMode),
					"-",
					strconv.Itoa(int(*vv.NumCacheNodes)),
				}
			} else if vv := v.ReplicationGroup; vv != nil {
				return []string{
					*vv.ReplicationGroupId,
					utils.TitleCase(*vv.Status),
					"Redis",
					v.ReplicationGroupEngineVersion,
					*vv.CacheNodeType,
					utils.TitleCase(string(vv.ClusterMode)),
					strconv.Itoa(len(vv.NodeGroups)),
					strconv.Itoa(len(vv.MemberClusters)),
				}
			}
			return []string{"", "", "", "", "", "", "", ""}
		},
		func(v *TableView[model.ElastiCacheCluster]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
					Description: "Service Update Status",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						var updateActionsView *TableView[model.ElastiCacheUpdateAction]
						if item.CacheCluster != nil {
							updateActionsView = NewElastiCacheUpdateActions(repo, app, []string{id}, []string{}, "")
						} else if item.ReplicationGroup != nil {
							updateActionsView = NewElastiCacheUpdateActions(repo, app, []string{}, []string{id}, "")
						} else {
							return
						}
						app.AddAndSwitch(updateActionsView)
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						var tagsView *Tags
						if item.CacheCluster != nil {
							if item.CacheCluster.ARN == nil {
								return
							}
							tagsView = NewTags(repo, v.GetService(), *item.CacheCluster.ARN, app)
						} else if item.ReplicationGroup != nil {
							if item.ReplicationGroup.ARN == nil {
								return
							}
							tagsView = NewTags(repo, v.GetService(), *item.ReplicationGroup.ARN, app)
						} else {
							return
						}
						app.AddAndSwitch(tagsView)
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Clusters", func(r Repos, app *Application) Component {
		return NewElastiCacheClusters(r.ElastiCache, app)
	})
}
