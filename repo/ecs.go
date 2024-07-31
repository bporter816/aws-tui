package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/bporter816/aws-tui/model"
)

type ECS struct {
	ecsClient *ecs.Client
}

func NewECS(ecsClient *ecs.Client) *ECS {
	return &ECS{
		ecsClient: ecsClient,
	}
}

// Internal function to get all cluster arns
func (e ECS) listClusterArns() ([]string, error) {
	pg := ecs.NewListClustersPaginator(
		e.ecsClient,
		&ecs.ListClustersInput{},
	)
	var arns []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []string{}, err
		}
		arns = append(arns, out.ClusterArns...)
	}
	return arns, nil
}

func (e ECS) ListClusters() ([]model.ECSCluster, error) {
	arns, err := e.listClusterArns()
	if err != nil {
		return []model.ECSCluster{}, err
	}
	var clusters []model.ECSCluster
	out, err := e.ecsClient.DescribeClusters(
		context.TODO(),
		&ecs.DescribeClustersInput{
			Clusters: arns,
		},
	)
	if err != nil {
		return []model.ECSCluster{}, err
	}
	for _, v := range out.Clusters {
		// TODO handle failures
		clusters = append(clusters, model.ECSCluster(v))
	}
	return clusters, nil
}

func (e ECS) ListTags(resourceId string) (model.Tags, error) {
	out, err := e.ecsClient.ListTagsForResource(
		context.TODO(),
		&ecs.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceId),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
