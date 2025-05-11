package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/bporter816/aws-tui/internal/model"
)

type EKS struct {
	eksClient *eks.Client
}

func NewEKS(eksClient *eks.Client) *EKS {
	return &EKS{
		eksClient: eksClient,
	}
}

func (e EKS) describeCluster(name string) (eksTypes.Cluster, error) {
	out, err := e.eksClient.DescribeCluster(
		context.TODO(),
		&eks.DescribeClusterInput{
			Name: aws.String(name),
		},
	)
	if err != nil {
		return eksTypes.Cluster{}, err
	}
	return *out.Cluster, nil
}

func (e EKS) ListClusters() ([]model.EKSCluster, error) {
	pg := eks.NewListClustersPaginator(
		e.eksClient,
		&eks.ListClustersInput{},
	)
	var clusters []model.EKSCluster
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.EKSCluster{}, err
		}
		for _, v := range out.Clusters {
			tmp := v
			cluster := model.EKSCluster{Name: &tmp}
			if details, err := e.describeCluster(v); err == nil {
				cluster = model.EKSCluster(details)
			}
			clusters = append(clusters, cluster)
		}
	}
	return clusters, nil
}

func (e EKS) ListTags(arn string) (model.Tags, error) {
	out, err := e.eksClient.ListTagsForResource(
		context.TODO(),
		&eks.ListTagsForResourceInput{
			ResourceArn: aws.String(arn),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for k, v := range out.Tags {
		tags = append(tags, model.Tag{Key: k, Value: v})
	}
	return tags, nil
}
