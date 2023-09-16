package repo

import (
	"context"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/bporter816/aws-tui/model"
)

type Cloudfront struct {
	cfClient *cf.Client
}

func NewCloudfront(cfClient *cf.Client) *Cloudfront {
	return &Cloudfront{
		cfClient: cfClient,
	}
}

func (c Cloudfront) ListDistributions() ([]model.CloudfrontDistribution, error) {
	pg := cf.NewListDistributionsPaginator(
		c.cfClient,
		&cf.ListDistributionsInput{},
	)
	var distributions []model.CloudfrontDistribution
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil || out.DistributionList == nil {
			return []model.CloudfrontDistribution{}, err
		}
		for _, v := range out.DistributionList.Items {
			distributions = append(distributions, model.CloudfrontDistribution(v))
		}
	}
	return distributions, nil
}

func (c Cloudfront) ListFunctions() ([]model.CloudfrontFunction, error) {
	// ListFunctions doesn't have a paginator
	var functions []model.CloudfrontFunction
	var marker *string
	for {
		out, err := c.cfClient.ListFunctions(
			context.TODO(),
			&cf.ListFunctionsInput{
				Marker: marker,
			},
		)
		if err != nil || out.FunctionList == nil {
			return []model.CloudfrontFunction{}, err
		}
		for _, v := range out.FunctionList.Items {
			functions = append(functions, model.CloudfrontFunction(v))
		}
		marker = out.FunctionList.NextMarker
		if marker == nil {
			break
		}
	}
	return functions, nil
}
