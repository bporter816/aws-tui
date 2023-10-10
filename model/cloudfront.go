package model

import (
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type (
	CloudFrontDistribution                    cfTypes.DistributionSummary
	CloudFrontDistributionOrigin              cfTypes.Origin
	CloudFrontDistributionCacheBehavior       cfTypes.CacheBehavior
	CloudFrontDistributionCustomErrorResponse cfTypes.CustomErrorResponse
	CloudFrontInvalidation                    cfTypes.InvalidationSummary
	CloudFrontInvalidationPath                string
	CloudFrontFunction                        cfTypes.FunctionSummary
)
