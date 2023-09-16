package model

import (
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type (
	CloudfrontDistribution cfTypes.DistributionSummary
	CloudfrontFunction     cfTypes.FunctionSummary
)
