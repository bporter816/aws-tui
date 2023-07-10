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
	ecClient          *ec.Client
	app               *Application
	clusters          []ecTypes.CacheCluster
	replicationGroups []ecTypes.ReplicationGroup
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
	var arn string
	if row <= len(e.clusters) {
		arn = *e.clusters[row-1].ARN
	} else {
		// TODO check bounds
		arn = *e.replicationGroups[row-1-len(e.clusters)].ARN
	}
	tagsView := NewElasticacheTags(e.ecClient, ElasticacheResourceTypeCluster, arn, name, e.app)
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

	clustersPg := ec.NewDescribeCacheClustersPaginator(
		e.ecClient,
		&ec.DescribeCacheClustersInput{},
	)
	e.clusters = make([]ecTypes.CacheCluster, 0)
	for clustersPg.HasMorePages() {
		out, err := clustersPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		e.clusters = append(e.clusters, out.CacheClusters...)
	}

	replicationGroupsPg := ec.NewDescribeReplicationGroupsPaginator(
		e.ecClient,
		&ec.DescribeReplicationGroupsInput{},
	)
	e.replicationGroups = make([]ecTypes.ReplicationGroup, 0)
	for replicationGroupsPg.HasMorePages() {
		out, err := replicationGroupsPg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		e.replicationGroups = append(e.replicationGroups, out.ReplicationGroups...)
	}

	caser := cases.Title(language.English)
	var data [][]string
	for _, v := range e.clusters {
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
	}
	for _, v := range e.replicationGroups {
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
	}
	e.SetData(data)
}
