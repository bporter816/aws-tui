package main

import (
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ElasticacheClusters struct {
	*ui.Table
	repo  *repo.Elasticache
	app   *Application
	model []model.ElasticacheCluster
}

func NewElasticacheClusters(repo *repo.Elasticache, app *Application) *ElasticacheClusters {
	e := &ElasticacheClusters{
		Table: ui.NewTable([]string{
			"ID",
			"STATUS",
			"ENGINE",
			"VERSION",
			"NODE TYPE",
			"CLUSTER MODE",
			"SHARDS",
			"NODES",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ElasticacheClusters) GetService() string {
	return "Elasticache"
}

func (e ElasticacheClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (e ElasticacheClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	var tagsView *ElasticacheTags
	if e.model[row-1].CacheCluster != nil {
		if e.model[row-1].CacheCluster.ARN == nil {
			return
		}
		tagsView = NewElasticacheTags(e.repo, ElasticacheResourceTypeCluster, *e.model[row-1].CacheCluster.ARN, name, e.app)
	} else if e.model[row-1].ReplicationGroup != nil {
		if e.model[row-1].ReplicationGroup.ARN == nil {
			return
		}
		tagsView = NewElasticacheTags(e.repo, ElasticacheResourceTypeCluster, *e.model[row-1].ReplicationGroup.ARN, name, e.app)
	}
	e.app.AddAndSwitch(tagsView)
}

func (e ElasticacheClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ElasticacheClusters) Render() {
	model, err := e.repo.ListClusters()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		if vv := v.CacheCluster; vv != nil {
			// skip clusters in replication groups as those are retrieved from DescribeReplicationGroups
			if vv.ReplicationGroupId != nil {
				continue
			}
			var clusterMode string = "-"
			if *vv.Engine == "redis" {
				clusterMode = string(ecTypes.ClusterModeDisabled)
			}
			data = append(data, []string{
				*vv.CacheClusterId,
				utils.TitleCase(*vv.CacheClusterStatus),
				utils.TitleCase(*vv.Engine),
				*vv.EngineVersion,
				*vv.CacheNodeType,
				utils.TitleCase(clusterMode),
				"-",
				strconv.Itoa(int(*vv.NumCacheNodes)),
			})
		} else if vv := v.ReplicationGroup; vv != nil {
			data = append(data, []string{
				*vv.ReplicationGroupId,
				utils.TitleCase(*vv.Status),
				"Redis",
				v.ReplicationGroupEngineVersion,
				*vv.CacheNodeType,
				utils.TitleCase(string(vv.ClusterMode)),
				strconv.Itoa(len(vv.NodeGroups)),
				strconv.Itoa(len(vv.MemberClusters)),
			})
		}
	}
	e.SetData(data)
}
