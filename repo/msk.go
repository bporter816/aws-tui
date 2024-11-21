package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	msk "github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/bporter816/aws-tui/model"
)

type MSK struct {
	mskClient *msk.Client
}

func NewMSK(mskClient *msk.Client) *MSK {
	return &MSK{
		mskClient: mskClient,
	}
}

func (m MSK) ListClusters() ([]model.MSKCluster, error) {
	pg := msk.NewListClustersV2Paginator(
		m.mskClient,
		&msk.ListClustersV2Input{},
	)
	var clusters []model.MSKCluster
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.MSKCluster{}, err
		}
		for _, v := range out.ClusterInfoList {
			clusters = append(clusters, model.MSKCluster(v))
		}
	}
	return clusters, nil
}

func (m MSK) ListTags(resourceId string) (model.Tags, error) {
	out, err := m.mskClient.ListTagsForResource(
		context.TODO(),
		&msk.ListTagsForResourceInput{
			ResourceArn: aws.String(resourceId),
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
