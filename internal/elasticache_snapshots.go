package internal

import (
	"strconv"
	"strings"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ElastiCacheSnapshots struct {
	*ui.Table
	view.ElastiCache
	repo  *repo.ElastiCache
	app   *Application
	model []model.ElastiCacheSnapshot
}

func NewElastiCacheSnapshots(repo *repo.ElastiCache, app *Application) *ElastiCacheSnapshots {
	e := &ElastiCacheSnapshots{
		Table: ui.NewTable([]string{
			"NAME",
			"CLUSTER",
			"TYPE",
			"CREATED",
			"STATUS",
			"SHARDS",
			"SIZE",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElastiCacheSnapshots) GetLabels() []string {
	return []string{"Snapshots"}
}

func (e ElastiCacheSnapshots) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheSnapshots) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElastiCacheSnapshots) Render() {
	model, err := e.repo.ListSnapshots()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
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
		for _, v := range v.NodeSnapshots {
			if v.CacheSize != nil {
				sizes = append(sizes, *v.CacheSize)
			}
		}
		size = strings.Join(sizes, ", ")
		data = append(data, []string{
			utils.DerefString(v.SnapshotName, ""),
			cluster,
			snapshotType,
			created,
			status,
			shards,
			size,
		})
	}
	e.SetData(data)
}
