package main

type IAMPolicyType string

const (
	IAMPolicyTypeManaged             IAMPolicyType = "Managed"
	IAMPolicyTypeInline              IAMPolicyType = "Inline"
	IAMPolicyTypePermissionsBoundary IAMPolicyType = "Permissions Boundary"
	IAMPolicyTypeAssumeRolePolicy    IAMPolicyType = "Assume Role Policy"
)
