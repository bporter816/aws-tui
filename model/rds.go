package model

import (
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type (
	RDSCluster  rdsTypes.DBCluster
	RDSInstance rdsTypes.DBInstance
)
