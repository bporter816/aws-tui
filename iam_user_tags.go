package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
)

type IAMUserTags struct {
	*ui.Table
	iamClient *iam.Client
	userName  string
	app       *Application
}

func NewIAMUserTags(iamClient *iam.Client, userName string, app *Application) *IAMUserTags {
	i := &IAMUserTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		iamClient: iamClient,
		userName:  userName,
		app:       app,
	}
	return i
}

func (i IAMUserTags) GetService() string {
	return "IAM"
}

func (i IAMUserTags) GetLabels() []string {
	return []string{i.userName, "Tags"}
}

func (i IAMUserTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMUserTags) Render() {
	pg := iam.NewListUserTagsPaginator(
		i.iamClient,
		&iam.ListUserTagsInput{
			UserName: aws.String(i.userName),
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
