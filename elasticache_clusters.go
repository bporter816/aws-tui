package main

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
)

type ElasticacheClusters struct {
	*Table
	ecClient *ec.Client
	app      *Application
	arns     []string
}

func NewElasticacheClusters(ecClient *ec.Client, app *Application) *ElasticacheClusters {
	e := &ElasticacheClusters{
		Table: NewTable([]string{
			"ID",
			"STATUS",
			"ENGINE",
			"VERSION",
			"NODE TYPE",
			"CLUSTER MODE",
			"SHARDS",
			"NODES",
		}, 1, 0),
		ecClient: ecClient,
		app:      app,
	}
	return e
}

func (e ElasticacheClusters) GetName() string {
	return "Elasticache | Clusters"
}

func (e ElasticacheClusters) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewElasticacheTags(e.ecClient, e.arns[row-1], e.app)
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

	caser := cases.Title(language.English)
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
			caser.String(*v.CacheClusterStatus),
			caser.String(*v.Engine),
			*v.EngineVersion,
			*v.CacheNodeType,
			caser.String(clusterMode),
			"-",
			strconv.Itoa(int(*v.NumCacheNodes)),
		})
		e.arns = append(e.arns, *v.ARN)
	}
	for _, v := range replicationGroups {
		firstMemberCluster := v.MemberClusters[0]
		data = append(data, []string{
			*v.ReplicationGroupId,
			caser.String(*v.Status),
			"Redis",
			clusterToEngineVersion[firstMemberCluster],
			*v.CacheNodeType,
			caser.String(string(v.ClusterMode)),
			strconv.Itoa(len(v.NodeGroups)),
			strconv.Itoa(len(v.MemberClusters)),
		})
		e.arns = append(e.arns, *v.ARN)
	}
	e.SetData(data)
}
