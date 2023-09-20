package repo

import (
	"context"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/bporter816/aws-tui/model"
)

type Elasticache struct {
	ecClient *ec.Client
}

func NewElasticache(ecClient *ec.Client) *Elasticache {
	return &Elasticache{
		ecClient: ecClient,
	}
}

func (e Elasticache) ListParameters(parameterGroupName string) ([]model.ElasticacheParameter, error) {
	pg := ec.NewDescribeCacheParametersPaginator(
		e.ecClient,
		&ec.DescribeCacheParametersInput{
			CacheParameterGroupName: aws.String(parameterGroupName),
		},
	)
	var parameters []model.ElasticacheParameter
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.ElasticacheParameter{}, err
		}
		for _, v := range out.Parameters {
			parameters = append(parameters, model.ElasticacheParameter(v))
		}
	}
	return parameters, nil
}

func (e Elasticache) ListTags(arn string) (model.Tags, error) {
	out, err := e.ecClient.ListTagsForResource(
		context.TODO(),
		&ec.ListTagsForResourceInput{
			ResourceName: aws.String(arn),
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
