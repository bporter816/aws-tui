package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/model"
)

type IAM struct {
	iamClient *iam.Client
}

func NewIAM(iamClient *iam.Client) *IAM {
	return &IAM{
		iamClient: iamClient,
	}
}

func (i IAM) getAccessKeyLastUsed(accessKeyId string) (iamTypes.AccessKeyLastUsed, error) {
	out, err := i.iamClient.GetAccessKeyLastUsed(
		context.TODO(),
		&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: aws.String(accessKeyId),
		},
	)
	if err != nil || out.AccessKeyLastUsed == nil {
		return iamTypes.AccessKeyLastUsed{}, err
	}
	return *out.AccessKeyLastUsed, nil
}

func (i IAM) ListAccessKeys(userName string) ([]model.IAMAccessKey, error) {
	pg := iam.NewListAccessKeysPaginator(
		i.iamClient,
		&iam.ListAccessKeysInput{
			UserName: aws.String(userName),
		},
	)
	var accessKeys []model.IAMAccessKey
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []model.IAMAccessKey{}, err
		}
		for _, v := range out.AccessKeyMetadata {
			m := model.IAMAccessKey{AccessKeyMetadata: v}
			if v.AccessKeyId != nil {
				lastUsed, err := i.getAccessKeyLastUsed(*v.AccessKeyId)
				if err == nil {
					m.LastUsed = lastUsed
				}
			}
			accessKeys = append(accessKeys, m)
		}
	}
	return accessKeys, nil
}

// TODO combine these functions and abstract it?
func (i IAM) ListRoleTags(roleName string) (model.Tags, error) {
	pg := iam.NewListRoleTagsPaginator(
		i.iamClient,
		&iam.ListRoleTagsInput{
			RoleName: aws.String(roleName),
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
		}

	}
	return tags, nil
}

func (i IAM) ListUserTags(userName string) (model.Tags, error) {
	pg := iam.NewListUserTagsPaginator(
		i.iamClient,
		&iam.ListUserTagsInput{
			UserName: aws.String(userName),
		},
	)
	var tags model.Tags
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return model.Tags{}, err
		}
		for _, v := range out.Tags {
			tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
		}
	}
	return tags, nil
}
