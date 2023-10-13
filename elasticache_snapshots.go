package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

type ElastiCacheSnapshots struct {
	*ui.Table
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

func (e ElastiCacheSnapshots) GetService() string {
	return "ElastiCache"
}

func (e ElastiCacheSnapshots) GetLabels() []string {
	return []string{"Snapshots"}
}

func (e ElastiCacheSnapshots) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if arn := e.model[row-1].ARN; arn != nil {
		tagsView := NewElastiCacheTags(e.repo, *arn, name, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ElastiCacheSnapshots) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		var name, cluster, snapshotType, created, status, size string
		shards := "-"
		if v.SnapshotName != nil {
			name = *v.SnapshotName
		}
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
			name,
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
