package main

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
)
