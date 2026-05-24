package repo

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/internal/model"
	"net/url"
	"strings"
)

type IAM struct {
	iamClient *iam.Client
}

func NewIAM(iamClient *iam.Client) *IAM {
	return &IAM{
		iamClient: iamClient,
	}
}

func (i IAM) ListAccountAliases(ctx context.Context) ([]string, error) {
	out, err := i.iamClient.ListAccountAliases(
		ctx,
		&iam.ListAccountAliasesInput{},
	)
	if err != nil {
		return []string{}, err
	}
	return out.AccountAliases, nil
}

func (i IAM) ListUsers(ctx context.Context, groupName *string) ([]model.IAMUser, error) {
	var users []model.IAMUser
	if groupName == nil {
		pg := iam.NewListUsersPaginator(
			i.iamClient,
			&iam.ListUsersInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(ctx)
			if err != nil {
				return []model.IAMUser{}, err
			}
			for _, v := range out.Users {
				users = append(users, model.IAMUser(v))
			}
		}
	} else {
		pg := iam.NewGetGroupPaginator(
			i.iamClient,
			&iam.GetGroupInput{
				GroupName: groupName,
			},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(ctx)
			if err != nil {
				return []model.IAMUser{}, err
			}
			for _, v := range out.Users {
				users = append(users, model.IAMUser(v))
			}
		}
	}
	return users, nil
}

func (i IAM) ListGroups(ctx context.Context, userName *string) ([]model.IAMGroup, error) {
	var groups []model.IAMGroup
	if userName == nil {
		pg := iam.NewListGroupsPaginator(
			i.iamClient,
			&iam.ListGroupsInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(ctx)
			if err != nil {
				return []model.IAMGroup{}, err
			}
			for _, v := range out.Groups {
				groups = append(groups, model.IAMGroup(v))
			}
		}
	} else {
		pg := iam.NewListGroupsForUserPaginator(
			i.iamClient,
			&iam.ListGroupsForUserInput{
				UserName: aws.String(*userName),
			},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(ctx)
			if err != nil {
				return []model.IAMGroup{}, err
			}
			for _, v := range out.Groups {
				groups = append(groups, model.IAMGroup(v))
			}
		}
	}
	return groups, nil
}

func (i IAM) ListRoles(ctx context.Context) ([]model.IAMRole, error) {
	pg := iam.NewListRolesPaginator(
		i.iamClient,
		&iam.ListRolesInput{},
	)
	var roles []model.IAMRole
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []model.IAMRole{}, err
		}
		for _, v := range out.Roles {
			roles = append(roles, model.IAMRole(v))
		}
	}
	return roles, nil
}

// TODO split this up
func (i IAM) ListPolicies(ctx context.Context, id *string, identityType model.IAMIdentityType) ([]model.IAMPolicy, error) {
	var policies []model.IAMPolicy
	// inline policies
	if id != nil {
		switch identityType {
		case model.IAMIdentityTypeUser:
			pg := iam.NewListUserPoliciesPaginator(
				i.iamClient,
				&iam.ListUserPoliciesInput{
					UserName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return []model.IAMPolicy{}, err
				}
				for _, v := range out.PolicyNames {
					policies = append(policies, model.IAMPolicy{
						Name:       v,
						PolicyType: model.IAMPolicyTypeInline,
					})
				}
			}
		case model.IAMIdentityTypeRole:
			pg := iam.NewListRolePoliciesPaginator(
				i.iamClient,
				&iam.ListRolePoliciesInput{
					RoleName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return []model.IAMPolicy{}, err
				}
				for _, v := range out.PolicyNames {
					policies = append(policies, model.IAMPolicy{
						Name:       v,
						PolicyType: model.IAMPolicyTypeInline,
					})
				}
			}
		case model.IAMIdentityTypeGroup:
			pg := iam.NewListGroupPoliciesPaginator(
				i.iamClient,
				&iam.ListGroupPoliciesInput{
					GroupName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return []model.IAMPolicy{}, err
				}
				for _, v := range out.PolicyNames {
					policies = append(policies, model.IAMPolicy{
						Name:       v,
						PolicyType: model.IAMPolicyTypeInline,
					})
				}
			}
		default:
			return []model.IAMPolicy{}, errors.New("invalid identity type for policy list")
		}
	}

	// managed policies
	if id == nil {
		// all policies in account
		pg := iam.NewListPoliciesPaginator(
			i.iamClient,
			&iam.ListPoliciesInput{},
		)
		for pg.HasMorePages() {
			out, err := pg.NextPage(ctx)
			if err != nil {
				return policies, err
			}
			for _, v := range out.Policies {
				var name string
				if v.PolicyName != nil {
					name = *v.PolicyName
				}
				policies = append(policies, model.IAMPolicy{
					Name:       name,
					PolicyType: model.IAMPolicyTypeManaged,
					Arn:        v.Arn,
				})
			}
		}
	} else {
		// policies attached to a user, role, or group
		switch identityType {
		case model.IAMIdentityTypeUser:
			pg := iam.NewListAttachedUserPoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedUserPoliciesInput{
					UserName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return policies, err
				}
				for _, v := range out.AttachedPolicies {
					var name string
					if v.PolicyName != nil {
						name = *v.PolicyName
					}
					policies = append(policies, model.IAMPolicy{
						Name:       name,
						PolicyType: model.IAMPolicyTypeManaged,
						Arn:        v.PolicyArn,
					})
				}
			}
		case model.IAMIdentityTypeRole:
			pg := iam.NewListAttachedRolePoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedRolePoliciesInput{
					RoleName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return policies, err
				}
				for _, v := range out.AttachedPolicies {
					var name string
					if v.PolicyName != nil {
						name = *v.PolicyName
					}
					policies = append(policies, model.IAMPolicy{
						Name:       name,
						PolicyType: model.IAMPolicyTypeManaged,
						Arn:        v.PolicyArn,
					})
				}
			}
		case model.IAMIdentityTypeGroup:
			pg := iam.NewListAttachedGroupPoliciesPaginator(
				i.iamClient,
				&iam.ListAttachedGroupPoliciesInput{
					GroupName: id,
				},
			)
			for pg.HasMorePages() {
				out, err := pg.NextPage(ctx)
				if err != nil {
					return policies, err
				}
				for _, v := range out.AttachedPolicies {
					var name string
					if v.PolicyName != nil {
						name = *v.PolicyName
					}
					policies = append(policies, model.IAMPolicy{
						Name:       name,
						PolicyType: model.IAMPolicyTypeManaged,
						Arn:        v.PolicyArn,
					})
				}
			}
		default:
			return policies, errors.New("invalid identity type for policy list")
		}
	}
	return policies, nil
}

func (i IAM) getIAMManagedPolicyCurrentVersion(ctx context.Context, policyArn string) (string, error) {
	// get the managed policy
	policyOut, err := i.iamClient.GetPolicy(
		ctx,
		&iam.GetPolicyInput{
			PolicyArn: aws.String(policyArn),
		},
	)
	if err != nil || policyOut.Policy == nil || policyOut.Policy.DefaultVersionId == nil {
		return "", err
	}

	// get the current version of the policy
	versionOut, err := i.iamClient.GetPolicyVersion(
		ctx,
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

func (i IAM) GetIAMManagedPolicy(ctx context.Context, policyArn string) (string, error) {
	return i.getIAMManagedPolicyCurrentVersion(ctx, policyArn)
}

func (i IAM) GetIAMInlinePolicy(ctx context.Context, identityType model.IAMIdentityType, identityName string, policyName string) (string, error) {
	var policyDocument *string
	switch identityType {
	case model.IAMIdentityTypeUser:
		out, err := i.iamClient.GetUserPolicy(
			ctx,
			&iam.GetUserPolicyInput{
				UserName:   aws.String(identityName),
				PolicyName: aws.String(policyName),
			},
		)
		if err != nil {
			return "", err
		}
		policyDocument = out.PolicyDocument
	case model.IAMIdentityTypeRole:
		out, err := i.iamClient.GetRolePolicy(
			ctx,
			&iam.GetRolePolicyInput{
				RoleName:   aws.String(identityName),
				PolicyName: aws.String(policyName),
			},
		)
		if err != nil {
			return "", err
		}
		policyDocument = out.PolicyDocument
	case model.IAMIdentityTypeGroup:
		out, err := i.iamClient.GetGroupPolicy(
			ctx,
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
		return "", errors.New("invalid identity type for inline policy, support types are user, role, and group")
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

func (i IAM) GetIAMPermissionsBoundary(ctx context.Context, name string, identityType model.IAMIdentityType) (string, error) {
	var attachment *iamTypes.AttachedPermissionsBoundary
	switch identityType {
	case model.IAMIdentityTypeUser:
		out, err := i.iamClient.GetUser(
			ctx,
			&iam.GetUserInput{
				UserName: aws.String(name),
			},
		)
		if err != nil || out.User == nil {
			return "", err
		}
		attachment = out.User.PermissionsBoundary
	case model.IAMIdentityTypeRole:
		out, err := i.iamClient.GetRole(
			ctx,
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

	return i.getIAMManagedPolicyCurrentVersion(ctx, *attachment.PermissionsBoundaryArn)
}

func (i IAM) GetIAMAssumeRolePolicy(ctx context.Context, roleName string) (string, error) {
	out, err := i.iamClient.GetRole(
		ctx,
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

func (i IAM) getAccessKeyLastUsed(ctx context.Context, accessKeyId string) (iamTypes.AccessKeyLastUsed, error) {
	out, err := i.iamClient.GetAccessKeyLastUsed(
		ctx,
		&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: aws.String(accessKeyId),
		},
	)
	if err != nil || out.AccessKeyLastUsed == nil {
		return iamTypes.AccessKeyLastUsed{}, err
	}
	return *out.AccessKeyLastUsed, nil
}

func (i IAM) ListAccessKeys(ctx context.Context, userName string) ([]model.IAMAccessKey, error) {
	pg := iam.NewListAccessKeysPaginator(
		i.iamClient,
		&iam.ListAccessKeysInput{
			UserName: aws.String(userName),
		},
	)
	var accessKeys []model.IAMAccessKey
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return []model.IAMAccessKey{}, err
		}
		for _, v := range out.AccessKeyMetadata {
			m := model.IAMAccessKey{AccessKeyMetadata: v}
			if v.AccessKeyId != nil {
				lastUsed, err := i.getAccessKeyLastUsed(ctx, *v.AccessKeyId)
				if err == nil {
					m.LastUsed = lastUsed
				}
			}
			accessKeys = append(accessKeys, m)
		}
	}
	return accessKeys, nil
}

func (i IAM) ListTags(ctx context.Context, typeAndName string) (model.Tags, error) {
	parts := strings.Split(typeAndName, ":")
	if len(parts) != 2 {
		return model.Tags{}, errors.New("must specify type and id for iam tags")
	}
	switch parts[0] {
	case "role":
		return i.listRoleTags(ctx, parts[1])
	case "user":
		return i.listUserTags(ctx, parts[1])
	default:
		return model.Tags{}, errors.New("must get iam tags for a role or a user")
	}
}

func (i IAM) listRoleTags(ctx context.Context, roleName string) (model.Tags, error) {
	pg := iam.NewListRoleTagsPaginator(
		i.iamClient,
		&iam.ListRoleTagsInput{
			RoleName: aws.String(roleName),
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
		}

	}
	return tags, nil
}

func (i IAM) listUserTags(ctx context.Context, userName string) (model.Tags, error) {
	pg := iam.NewListUserTagsPaginator(
		i.iamClient,
		&iam.ListUserTagsInput{
			UserName: aws.String(userName),
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(ctx)
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
		}
	}
	return tags, nil
}
