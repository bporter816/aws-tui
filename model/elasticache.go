package model

import (
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

type (
	// A cluster is either a cache cluster or a replication group
	ElasticacheCluster struct {
		CacheCluster *ecTypes.CacheCluster
		ReplicationGroup *ecTypes.ReplicationGroup
		ReplicationGroupEngineVersion string
	}
	ElasticacheEvent ecTypes.Event
	ElasticacheReservedNode ecTypes.ReservedCacheNode
	ElasticacheSnapshot ecTypes.Snapshot
	ElasticacheParameterGroup ecTypes.CacheParameterGroup
	ElasticacheParameter ecTypes.Parameter
	ElasticacheSubnetGroup ecTypes.CacheSubnetGroup
	ElasticacheUser ecTypes.User
	ElasticacheGroup ecTypes.UserGroup
)
