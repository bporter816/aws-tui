package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/bporter816/aws-tui/model"
)

type SSM struct {
	ssmClient *ssm.Client
}

func NewSSM(ssmClient *ssm.Client) *SSM {
	return &SSM{
		ssmClient: ssmClient,
	}
}

func (s SSM) ListParameters() ([]model.SSMParameter, error) {
	pg := ssm.NewDescribeParametersPaginator(
		s.ssmClient,
		&ssm.DescribeParametersInput{},
	)
	var parameters []model.SSMParameter
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.SSMParameter{}, err
		}
		for _, v := range out.Parameters {
			parameters = append(parameters, model.SSMParameter(v))
		}
	}
	return parameters, nil
}

func (s SSM) ListTags(resourceId string) (model.Tags, error) {
	parts := strings.Split(resourceId, ":")
	if len(parts) != 2 {
		return model.Tags{}, errors.New("must provide type and arn for ssm tags")
	}
	var resourceType ssmTypes.ResourceTypeForTagging
	switch parts[0] {
	case "parameter":
		resourceType = ssmTypes.ResourceTypeForTaggingParameter
	default:
		return model.Tags{}, errors.New("unknown type for ssm tags")
	}

	out, err := s.ssmClient.ListTagsForResource(
		context.TODO(),
		&ssm.ListTagsForResourceInput{
			ResourceId:   aws.String(parts[1]),
			ResourceType: resourceType,
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
