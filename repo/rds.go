package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/model"
)

type RDS struct {
	rdsClient *rds.Client
}

func NewRDS(rdsClient *rds.Client) *RDS {
	return &RDS{
		rdsClient: rdsClient,
	}
}

func (r RDS) ListClusters(filters []rdsTypes.Filter) ([]model.RDSCluster, error) {
	pg := rds.NewDescribeDBClustersPaginator(
		r.rdsClient,
		&rds.DescribeDBClustersInput{
			Filters: filters,
		},
	)
	var clusters []model.RDSCluster
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.RDSCluster{}, err
		}
		for _, v := range out.DBClusters {
			clusters = append(clusters, model.RDSCluster(v))
		}
	}
	return clusters, nil
}

func (r RDS) ListGlobalClusters() ([]model.RDSGlobalCluster, error) {
	pg := rds.NewDescribeGlobalClustersPaginator(
		r.rdsClient,
		&rds.DescribeGlobalClustersInput{},
	)
	var globalClusters []model.RDSGlobalCluster
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.RDSGlobalCluster{}, err
		}
		for _, v := range out.GlobalClusters {
			globalClusters = append(globalClusters, model.RDSGlobalCluster(v))
		}
	}
	return globalClusters, nil
}

func (r RDS) ListInstances(filters []rdsTypes.Filter) ([]model.RDSInstance, error) {
	pg := rds.NewDescribeDBInstancesPaginator(
		r.rdsClient,
		&rds.DescribeDBInstancesInput{
			Filters: filters,
		},
	)
	var instances []model.RDSInstance
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.RDSInstance{}, err
		}
		for _, v := range out.DBInstances {
			instances = append(instances, model.RDSInstance(v))
		}
	}
	return instances, nil
}

func (r RDS) ListTags(resourceId string) (model.Tags, error) {
	out, err := r.rdsClient.ListTagsForResource(
		context.TODO(),
		&rds.ListTagsForResourceInput{
			ResourceName: aws.String(resourceId),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.TagList {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
