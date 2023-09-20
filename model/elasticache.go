package model

import (
	ecTypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
)

type (
	ElasticacheEvent ecTypes.Event
	ElasticacheReservedNode ecTypes.ReservedCacheNode
	ElasticacheSnapshot ecTypes.Snapshot
	ElasticacheParameterGroup ecTypes.CacheParameterGroup
	ElasticacheParameter ecTypes.Parameter
)
