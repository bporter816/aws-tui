package model

import (
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type (
	CloudfrontDistribution                    cfTypes.DistributionSummary
	CloudfrontDistributionOrigin              cfTypes.Origin
	CloudfrontDistributionCacheBehavior       cfTypes.CacheBehavior
	CloudfrontDistributionCustomErrorResponse cfTypes.CustomErrorResponse
	CloudfrontFunction                        cfTypes.FunctionSummary
)
