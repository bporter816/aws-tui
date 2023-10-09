package model

import (
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type (
	IAMUser   iamTypes.User
	IAMGroup  iamTypes.Group
	IAMRole   iamTypes.Role
	IAMPolicy struct {
		Name       string
		PolicyType IAMPolicyType
		Arn        *string // ARN for managed policies
	}
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

type IAMPolicyType string

const (
	IAMPolicyTypeManaged             IAMPolicyType = "Managed"
	IAMPolicyTypeInline              IAMPolicyType = "Inline"
	IAMPolicyTypePermissionsBoundary IAMPolicyType = "Permissions Boundary"
	IAMPolicyTypeAssumeRolePolicy    IAMPolicyType = "Assume Role Policy"
)
