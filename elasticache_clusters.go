package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type ElasticacheClusters struct {
	*ui.Table
	repo     *repo.Elasticache
	ecClient *ec.Client
	app      *Application
	arns     []string
}

func NewElasticacheClusters(repo *repo.Elasticache, ecClient *ec.Client, app *Application) *ElasticacheClusters {
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
		repo:     repo,
		ecClient: ecClient,
		app:      app,
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
	tagsView := NewElasticacheTags(e.repo, ElasticacheResourceTypeCluster, e.arns[row-1], name, e.app)
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
	// DescribeReplicationGroups doesn't return engine version, so we have to get it from the list of member cluster names
	clusterToEngineVersion := make(map[string]string)
	e.arns = make([]string, 0)

	clustersPg := ec.NewDescribeCacheClustersPaginator(
		e.ecClient,
		&ec.DescribeCacheClustersInput{},
	)
	var clusters []ecTypes.CacheCluster
	for clustersPg.HasMorePages() {
		out, err := clustersPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		clusters = append(clusters, out.CacheClusters...)
	}

	replicationGroupsPg := ec.NewDescribeReplicationGroupsPaginator(
		e.ecClient,
		&ec.DescribeReplicationGroupsInput{},
	)
	var replicationGroups []ecTypes.ReplicationGroup
	for replicationGroupsPg.HasMorePages() {
		out, err := replicationGroupsPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		replicationGroups = append(replicationGroups, out.ReplicationGroups...)
	}

	var data [][]string
	for _, v := range clusters {
		// skip clusters in replication groups as those are retrieved from DescribeReplicationGroups
		if v.ReplicationGroupId != nil {
			clusterToEngineVersion[*v.CacheClusterId] = *v.EngineVersion
			continue
		}
		var clusterMode string = "-"
		if *v.Engine == "redis" {
			clusterMode = string(ecTypes.ClusterModeDisabled)
		}
		data = append(data, []string{
			*v.CacheClusterId,
			utils.TitleCase(*v.CacheClusterStatus),
			utils.TitleCase(*v.Engine),
			*v.EngineVersion,
			*v.CacheNodeType,
			utils.TitleCase(clusterMode),
			"-",
			strconv.Itoa(int(*v.NumCacheNodes)),
		})
		e.arns = append(e.arns, *v.ARN)
	}
	for _, v := range replicationGroups {
		firstMemberCluster := v.MemberClusters[0]
		data = append(data, []string{
			*v.ReplicationGroupId,
			utils.TitleCase(*v.Status),
			"Redis",
			clusterToEngineVersion[firstMemberCluster],
			*v.CacheNodeType,
			utils.TitleCase(string(v.ClusterMode)),
			strconv.Itoa(len(v.NodeGroups)),
			strconv.Itoa(len(v.MemberClusters)),
		})
		e.arns = append(e.arns, *v.ARN)
	}
	e.SetData(data)
}
