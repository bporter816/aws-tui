package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/bporter816/aws-tui/model"
)

type CloudFront struct {
	cfClient *cf.Client
}

func NewCloudFront(cfClient *cf.Client) *CloudFront {
	return &CloudFront{
		cfClient: cfClient,
	}
}

func (c CloudFront) getDistributionConfig(distributionId string) (*cfTypes.DistributionConfig, error) {
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

func (c CloudFront) ListDistributions() ([]model.CloudFrontDistribution, error) {
	pg := cf.NewListDistributionsPaginator(
		c.cfClient,
		&cf.ListDistributionsInput{},
	)
	var distributions []model.CloudFrontDistribution
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil || out.DistributionList == nil {
			return []model.CloudFrontDistribution{}, err
		}
		for _, v := range out.DistributionList.Items {
			distributions = append(distributions, model.CloudFrontDistribution(v))
		}
	}
	return distributions, nil
}

func (c CloudFront) GetDistributionOrigins(distributionId string) ([]model.CloudFrontDistributionOrigin, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil || out.Origins == nil {
		return []model.CloudFrontDistributionOrigin{}, err
	}
	var origins []model.CloudFrontDistributionOrigin
	for _, v := range out.Origins.Items {
		origins = append(origins, model.CloudFrontDistributionOrigin(v))
	}
	return origins, nil
}

func (c CloudFront) GetDistributionCacheBehaviors(distributionId string) ([]model.CloudFrontDistributionCacheBehavior, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil {
		return []model.CloudFrontDistributionCacheBehavior{}, err
	}
	var cacheBehaviors []model.CloudFrontDistributionCacheBehavior
	if out.CacheBehaviors != nil {
		for _, v := range out.CacheBehaviors.Items {
			cacheBehaviors = append(cacheBehaviors, model.CloudFrontDistributionCacheBehavior(v))
		}
	}
	// the default cache behavior is a different type, so merge them
	if d := out.DefaultCacheBehavior; d != nil {
		cacheBehaviors = append(cacheBehaviors, model.CloudFrontDistributionCacheBehavior{
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

func (c CloudFront) GetDistributionCustomErrorResponses(distributionId string) ([]model.CloudFrontDistributionCustomErrorResponse, error) {
	out, err := c.getDistributionConfig(distributionId)
	if err != nil || out.CustomErrorResponses == nil {
		return []model.CloudFrontDistributionCustomErrorResponse{}, err
	}
	var customErrorResponses []model.CloudFrontDistributionCustomErrorResponse
	for _, v := range out.CustomErrorResponses.Items {
		customErrorResponses = append(customErrorResponses, model.CloudFrontDistributionCustomErrorResponse(v))
	}
	return customErrorResponses, nil
}

func (c CloudFront) ListInvalidations(distributionId string) ([]model.CloudFrontInvalidation, error) {
	pg := cf.NewListInvalidationsPaginator(
		c.cfClient,
		&cf.ListInvalidationsInput{
			DistributionId: aws.String(distributionId),
		},
	)
	var invalidations []model.CloudFrontInvalidation
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil || out.InvalidationList == nil {
			return []model.CloudFrontInvalidation{}, err
		}
		for _, v := range out.InvalidationList.Items {
			invalidations = append(invalidations, model.CloudFrontInvalidation(v))
		}
	}
	return invalidations, nil
}

func (c CloudFront) ListInvalidationPaths(distributionId string, invalidationId string) ([]model.CloudFrontInvalidationPath, error) {
	out, err := c.cfClient.GetInvalidation(
		context.TODO(),
		&cf.GetInvalidationInput{
			DistributionId: aws.String(distributionId),
			Id:             aws.String(invalidationId),
		},
	)
	if err != nil || out.Invalidation == nil || out.Invalidation.InvalidationBatch == nil || out.Invalidation.InvalidationBatch.Paths == nil {
		return []model.CloudFrontInvalidationPath{}, err
	}
	var paths []model.CloudFrontInvalidationPath
	for _, v := range out.Invalidation.InvalidationBatch.Paths.Items {
		paths = append(paths, model.CloudFrontInvalidationPath(v))
	}
	return paths, nil
}

func (c CloudFront) ListFunctions() ([]model.CloudFrontFunction, error) {
	// ListFunctions doesn't have a paginator
	var functions []model.CloudFrontFunction
	var marker *string
	for {
		out, err := c.cfClient.ListFunctions(
			context.TODO(),
			&cf.ListFunctionsInput{
				Marker: marker,
			},
		)
		if err != nil || out.FunctionList == nil {
			return []model.CloudFrontFunction{}, err
		}
		for _, v := range out.FunctionList.Items {
			functions = append(functions, model.CloudFrontFunction(v))
		}
		marker = out.FunctionList.NextMarker
		if marker == nil {
			break
		}
	}
	return functions, nil
}

func (c CloudFront) ListTags(resourceId string) (model.Tags, error) {
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
