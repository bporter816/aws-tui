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
			PathPattern:             aws.String("Default (*)"),
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

func (c Cloudfront) ListInvalidations(distributionId string) ([]model.CloudfrontInvalidation, error) {
	pg := cf.NewListInvalidationsPaginator(
		c.cfClient,
		&cf.ListInvalidationsInput{
			DistributionId: aws.String(distributionId),
		},
	)
	var invalidations []model.CloudfrontInvalidation
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil || out.InvalidationList == nil {
			return []model.CloudfrontInvalidation{}, err
		}
		for _, v := range out.InvalidationList.Items {
			invalidations = append(invalidations, model.CloudfrontInvalidation(v))
		}
	}
	return invalidations, nil
}

func (c Cloudfront) ListInvalidationPaths(distributionId string, invalidationId string) ([]model.CloudfrontInvalidationPath, error) {
	out, err := c.cfClient.GetInvalidation(
		context.TODO(),
		&cf.GetInvalidationInput{
			DistributionId: aws.String(distributionId),
			Id:             aws.String(invalidationId),
		},
	)
	if err != nil || out.Invalidation == nil || out.Invalidation.InvalidationBatch == nil || out.Invalidation.InvalidationBatch.Paths == nil {
		return []model.CloudfrontInvalidationPath{}, err
	}
	var paths []model.CloudfrontInvalidationPath
	for _, v := range out.Invalidation.InvalidationBatch.Paths.Items {
		paths = append(paths, model.CloudfrontInvalidationPath(v))
	}
	return paths, nil
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

func (c Cloudfront) ListTags(resourceId string) (model.Tags, error) {
	out, err := c.cfClient.ListTagsForResource(
		context.TODO(),
		&cf.ListTagsForResourceInput{
			Resource: aws.String(resourceId),
		},
	)
	if err != nil || out.Tags == nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags.Items {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
