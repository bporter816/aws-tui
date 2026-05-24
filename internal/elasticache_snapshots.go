package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewElastiCacheSnapshots(repo *repo.ElastiCache, app *Application) *TableView[model.ElastiCacheSnapshot] {
	return NewTableView(
		"ElastiCache",
		[]string{"Snapshots"},
		[]string{"NAME", "CLUSTER", "TYPE", "CREATED", "STATUS", "SHARDS", "SIZE"},
		repo.ListSnapshots,
		func(v model.ElastiCacheSnapshot) []string {
			var cluster, snapshotType, created, status, size string
			shards := "-"
			if v.ReplicationGroupId != nil {
				cluster = *v.ReplicationGroupId
				shards = strconv.Itoa(len(v.NodeSnapshots))
			} else if v.CacheClusterId != nil {
				cluster = *v.CacheClusterId
			}
			if v.SnapshotSource != nil {
				snapshotType = utils.TitleCase(*v.SnapshotSource)
			}
			if len(v.NodeSnapshots) > 0 && v.NodeSnapshots[0].SnapshotCreateTime != nil {
				created = v.NodeSnapshots[0].SnapshotCreateTime.Format(utils.DefaultTimeFormat)
			}
			if v.SnapshotStatus != nil {
				status = utils.TitleCase(*v.SnapshotStatus)
			}
			// TODO do math and sum these up? don't know what units they could be
			sizes := make([]string, 0)
			for _, n := range v.NodeSnapshots {
				if n.CacheSize != nil {
					sizes = append(sizes, *n.CacheSize)
				}
			}
			size = strings.Join(sizes, ", ")
			return []string{
				utils.DerefString(v.SnapshotName, ""),
				cluster,
				snapshotType,
				created,
				status,
				shards,
				size,
			}
		},
		func(v *TableView[model.ElastiCacheSnapshot]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ARN == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ARN, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ElastiCache.Snapshots", func(r Repos, app *Application) Component {
		return NewElastiCacheSnapshots(r.ElastiCache, app)
	})
}
