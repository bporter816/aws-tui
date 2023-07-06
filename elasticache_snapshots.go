package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

type ElasticacheSnapshots struct {
	*Table
	ecClient *ec.Client
	app      *Application
}

func NewElasticacheSnapshots(ecClient *ec.Client, app *Application) *ElasticacheSnapshots {
	e := &ElasticacheSnapshots{
		Table: NewTable([]string{
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

func (e ElasticacheSnapshots) GetName() string {
	return "Elasticache | Snapshots"
}

func (e ElasticacheSnapshots) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheSnapshots) Render() {
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

	caser := cases.Title(language.English)
	var data [][]string
	for _, v := range snapshots {
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
			snapshotType = caser.String(*v.SnapshotSource)
		}
		if len(v.NodeSnapshots) > 0 && v.NodeSnapshots[0].SnapshotCreateTime != nil {
			created = v.NodeSnapshots[0].SnapshotCreateTime.Format("2006-01-02 15:04:05")
		}
		if v.SnapshotStatus != nil {
			status = caser.String(*v.SnapshotStatus)
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
