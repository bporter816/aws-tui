package model

import (
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type (
	RDSCluster                rdsTypes.DBCluster
	RDSGlobalCluster          rdsTypes.GlobalCluster
	RDSInstance               rdsTypes.DBInstance
	RDSClusterParameterGroup  rdsTypes.DBClusterParameterGroup
	RDSInstanceParameterGroup rdsTypes.DBParameterGroup
	RDSParameter              rdsTypes.Parameter
	RDSSubnetGroup            rdsTypes.DBSubnetGroup
)

func (r RDSClusterParameterGroup) Arn() string {
	return *r.DBClusterParameterGroupArn
}

func (r RDSInstanceParameterGroup) Arn() string {
	return *r.DBParameterGroupArn
}

type RDSParameterGroupType string

const RDSParameterGroupTypeCluster RDSParameterGroupType = "Cluster"
const RDSParameterGroupTypeInstance RDSParameterGroupType = "Instance"
