package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

type ElasticacheSnapshots struct {
	*ui.Table
	ecClient *ec.Client
	app      *Application
	arns     []string
}

func NewElasticacheSnapshots(ecClient *ec.Client, app *Application) *ElasticacheSnapshots {
	e := &ElasticacheSnapshots{
		Table: ui.NewTable([]string{
			"NAME",
			"CLUSTER",
			"TYPE",
			"CREATED",
			"STATUS",
			"SHARDS",
			"SIZE",
		}, 1, 0),
		ecClient: ecClient,
		app:      app,
	}
	return e
}

func (e ElasticacheSnapshots) GetService() string {
	return "Elasticache"
}

func (e ElasticacheSnapshots) GetLabels() []string {
	return []string{"Snapshots"}
}

func (e ElasticacheSnapshots) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewElasticacheTags(e.ecClient, ElasticacheResourceTypeSnapshot, e.arns[row-1], name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ElasticacheSnapshots) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElasticacheSnapshots) Render() {
	pg := ec.NewDescribeSnapshotsPaginator(
		e.ecClient,
		&ec.DescribeSnapshotsInput{},
	)
	var snapshots []ecTypes.Snapshot
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		snapshots = append(snapshots, out.Snapshots...)
	}

	var data [][]string
	e.arns = make([]string, len(snapshots))
	for i, v := range snapshots {
		e.arns[i] = *v.ARN
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
