package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
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
