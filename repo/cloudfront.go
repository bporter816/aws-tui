package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
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

func (c Cloudfront) getDistributionConfig(distributionId string) (*cfTypes.DistributionConfig, error) {
	out, err := c.cfClient.GetDistributionConfig(
		context.TODO(),
		&cf.GetDistributionConfigInput{
			Id: aws.String(distributionId),
		},
	)
	if err != nil || out.DistributionConfig == nil {
		return &cfTypes.DistributionConfig{}, err
	}
	return out.DistributionConfig, nil
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

func (c Cloudfront) GetDistributionOrigins(distributionId string) ([]model.CloudfrontDistributionOrigin, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil || out.Origins == nil {
		return []model.CloudfrontDistributionOrigin{}, err
	}
	var origins []model.CloudfrontDistributionOrigin
	for _, v := range out.Origins.Items {
		origins = append(origins, model.CloudfrontDistributionOrigin(v))
	}
	return origins, nil
}

func (c Cloudfront) GetDistributionCacheBehaviors(distributionId string) ([]model.CloudfrontDistributionCacheBehavior, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil {
		return []model.CloudfrontDistributionCacheBehavior{}, err
	}
	var cacheBehaviors []model.CloudfrontDistributionCacheBehavior
	if out.CacheBehaviors != nil {
		for _, v := range out.CacheBehaviors.Items {
			cacheBehaviors = append(cacheBehaviors, model.CloudfrontDistributionCacheBehavior(v))
		}
	}
	// the default cache behavior is a different type, so merge them
	if d := out.DefaultCacheBehavior; d != nil {
		cacheBehaviors = append(cacheBehaviors, model.CloudfrontDistributionCacheBehavior{
			PathPattern: aws.String("Default (*)"),
			TargetOriginId:          d.TargetOriginId,
			ViewerProtocolPolicy:    d.ViewerProtocolPolicy,
			CachePolicyId:           d.CachePolicyId,
			OriginRequestPolicyId:   d.OriginRequestPolicyId,
			ResponseHeadersPolicyId: d.ResponseHeadersPolicyId,
		})
	}
	return cacheBehaviors, nil
}

func (c Cloudfront) GetDistributionCustomErrorResponses(distributionId string) ([]model.CloudfrontDistributionCustomErrorResponse, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil || out.CustomErrorResponses == nil {
		return []model.CloudfrontDistributionCustomErrorResponse{}, err
	}
	var customErrorResponses []model.CloudfrontDistributionCustomErrorResponse
	for _, v := range out.CustomErrorResponses.Items {
		customErrorResponses = append(customErrorResponses, model.CloudfrontDistributionCustomErrorResponse(v))
	}
	return customErrorResponses, nil
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
