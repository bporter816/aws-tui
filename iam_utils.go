package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"net/url"
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

func getIAMManagedPolicyCurrentVersion(iamClient *iam.Client, policyArn string) (string, error) {
	// get the managed policy
	policyOut, err := iamClient.GetPolicy(
		context.TODO(),
		&iam.GetPolicyInput{
			PolicyArn: aws.String(policyArn),
		},
	)
	if err != nil || policyOut.Policy == nil || policyOut.Policy.DefaultVersionId == nil {
		return "", err
	}

	// get the current version of the policy
	versionOut, err := iamClient.GetPolicyVersion(
		context.TODO(),
		&iam.GetPolicyVersionInput{
			PolicyArn: aws.String(policyArn),
			VersionId: policyOut.Policy.DefaultVersionId, // TODO use aws.String?
		},
	)
	if err != nil || versionOut.PolicyVersion == nil || versionOut.PolicyVersion.Document == nil {
		return "", err
	}

	// decode the policy
	decodedStr, err := url.QueryUnescape(*versionOut.PolicyVersion.Document)
	if err != nil {
		return "", errors.New("error decoding policy document")
	}

	return decodedStr, nil
}

func getIAMManagedPolicy(iamClient *iam.Client, policyArn string) (string, error) {
	return getIAMManagedPolicyCurrentVersion(iamClient, policyArn)
}

func getIAMInlinePolicy(iamClient *iam.Client, identityType IAMIdentityType, identityName string, policyName string) (string, error) {
	var policyDocument *string
	switch identityType {
	case IAMIdentityTypeUser:
		out, err := iamClient.GetUserPolicy(
			context.TODO(),
			&iam.GetUserPolicyInput{
				UserName:   aws.String(identityName),
				PolicyName: aws.String(policyName),
			},
		)
		if err != nil {
			return "", err
		}
		policyDocument = out.PolicyDocument
	case IAMIdentityTypeRole:
		out, err := iamClient.GetRolePolicy(
			context.TODO(),
			&iam.GetRolePolicyInput{
				RoleName:   aws.String(identityName),
				PolicyName: aws.String(policyName),
			},
		)
		if err != nil {
			return "", err
		}
		policyDocument = out.PolicyDocument
	case IAMIdentityTypeGroup:
		out, err := iamClient.GetGroupPolicy(
			context.TODO(),
			&iam.GetGroupPolicyInput{
				GroupName:  aws.String(identityName),
				PolicyName: aws.String(policyName),
			},
		)
		if err != nil {
			return "", err
		}
		policyDocument = out.PolicyDocument
	default:
		panic("invalid identity type for inline policy, support types are user, role, and group")
	}
	if policyDocument == nil || *policyDocument == "" {
		return "", nil
	}

	decodedStr, err := url.QueryUnescape(*policyDocument)
	if err != nil {
		return "", errors.New("error decoding policy document")
	}

	return decodedStr, nil
}

func getIAMPermissionsBoundary(iamClient *iam.Client, name string, identityType IAMIdentityType) (string, error) {
	var attachment *iamTypes.AttachedPermissionsBoundary
	switch identityType {
	case IAMIdentityTypeUser:
		out, err := iamClient.GetUser(
			context.TODO(),
			&iam.GetUserInput{
				UserName: aws.String(name),
			},
		)
		if err != nil || out.User == nil {
			return "", err
		}
		attachment = out.User.PermissionsBoundary
	case IAMIdentityTypeRole:
		out, err := iamClient.GetRole(
			context.TODO(),
			&iam.GetRoleInput{
				RoleName: aws.String(name),
			},
		)
		if err != nil || out.Role == nil {
			return "", err
		}
		attachment = out.Role.PermissionsBoundary
	default:
		panic("invalid identity type for permissions boundary, supported types are user and role")
	}
	if attachment == nil || attachment.PermissionsBoundaryArn == nil || *attachment.PermissionsBoundaryArn == "" {
		return "", nil
	}

	return getIAMManagedPolicyCurrentVersion(iamClient, *attachment.PermissionsBoundaryArn)
}

func getIAMAssumeRolePolicy(iamClient *iam.Client, roleName string) (string, error) {
	out, err := iamClient.GetRole(
		context.TODO(),
		&iam.GetRoleInput{
			RoleName: aws.String(roleName),
		},
	)
	if err != nil || out.Role == nil || out.Role.AssumeRolePolicyDocument == nil {
		return "", err
	}

	decodedStr, err := url.QueryUnescape(*out.Role.AssumeRolePolicyDocument)
	if err != nil {
		return "", errors.New("error decoding policy document")
	}

	return decodedStr, nil
}
