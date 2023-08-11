package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
)

type IAMRoleTags struct {
	*ui.Table
	iamClient *iam.Client
	roleName  string
	app       *Application
}

func NewIAMRoleTags(iamClient *iam.Client, roleName string, app *Application) *IAMRoleTags {
	i := &IAMRoleTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		iamClient: iamClient,
		roleName:  roleName,
		app:       app,
	}
	return i
}

func (i IAMRoleTags) GetService() string {
	return "IAM"
}

func (i IAMRoleTags) GetLabels() []string {
	return []string{i.roleName, "Tags"}
}

func (i IAMRoleTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMRoleTags) Render() {
	pg := iam.NewListRoleTagsPaginator(
		i.iamClient,
		&iam.ListRoleTagsInput{
			RoleName: aws.String(i.roleName),
		},
	)
	var tags []iamTypes.Tag
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		tags = append(tags, out.Tags...)
	}

	var data [][]string
	for _, v := range tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	i.SetData(data)
}
