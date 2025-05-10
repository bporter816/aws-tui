package model

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type (
	ECSCluster ecsTypes.Cluster
	ECSService ecsTypes.Service
	ECSTask    ecsTypes.Task
)
