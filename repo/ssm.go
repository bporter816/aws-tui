package repo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
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
