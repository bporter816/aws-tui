package model

import (
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type (
	IAMUser iamTypes.User
	IAMGroup iamTypes.Group
	IAMRole iamTypes.Role
	IAMAccessKey struct {
		iamTypes.AccessKeyMetadata
		LastUsed iamTypes.AccessKeyLastUsed
	}
)

type IAMIdentityType string

const (
	IAMIdentityTypeUser  IAMIdentityType = "User"
	IAMIdentityTypeRole  IAMIdentityType = "Role"
	IAMIdentityTypeGroup IAMIdentityType = "Group"
	IAMIdentityTypeAll   IAMIdentityType = "All"
)
