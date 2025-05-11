package model

import (
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

type (
	// A cluster is either a cache cluster or a replication group
	ElastiCacheCluster struct {
		CacheCluster                  *ecTypes.CacheCluster
		ReplicationGroup              *ecTypes.ReplicationGroup
		ReplicationGroupEngineVersion string
	}
	ElastiCacheEvent          ecTypes.Event
	ElastiCacheReservedNode   ecTypes.ReservedCacheNode
	ElastiCacheSnapshot       ecTypes.Snapshot
	ElastiCacheParameterGroup ecTypes.CacheParameterGroup
	ElastiCacheParameter      ecTypes.Parameter
	ElastiCacheSubnetGroup    ecTypes.CacheSubnetGroup
	ElastiCacheUser           ecTypes.User
	ElastiCacheGroup          ecTypes.UserGroup
	ElastiCacheServiceUpdate  ecTypes.ServiceUpdate
	ElastiCacheUpdateAction   ecTypes.UpdateAction
)
