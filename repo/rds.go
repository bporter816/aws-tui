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
		&rds.DescribeDBClustersInput{},
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
