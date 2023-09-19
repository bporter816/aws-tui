package model

import (
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type (
	IAMAccessKey struct {
		iamTypes.AccessKeyMetadata
		LastUsed iamTypes.AccessKeyLastUsed
	}
)
