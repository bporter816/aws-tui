package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type IAMAccessKeys struct {
	*ui.Table
	iamClient *iam.Client
	userName  string
	app       *Application
}

func NewIAMAccessKeys(iamClient *iam.Client, userName string, app *Application) *IAMAccessKeys {
	i := &IAMAccessKeys{
		Table: ui.NewTable([]string{
			"ID",
			"CREATED",
			"STATUS",
			"LAST USED DATE",
			"LAST USED REGION",
			"LAST USED SERVICE",
		}, 1, 0),
		iamClient: iamClient,
		userName:  userName,
		app:       app,
	}
	return i
}

func (i IAMAccessKeys) GetService() string {
	return "IAM"
}

func (i IAMAccessKeys) GetLabels() []string {
	return []string{i.userName, "Access Keys"}
}

func (i IAMAccessKeys) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMAccessKeys) Render() {
	var accessKeys []iamTypes.AccessKeyMetadata
	pg := iam.NewListAccessKeysPaginator(
		i.iamClient,
		&iam.ListAccessKeysInput{
			UserName: aws.String(i.userName),
		},
	)
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		accessKeys = append(accessKeys, out.AccessKeyMetadata...)
	}

	var data [][]string

	for _, v := range accessKeys {
		var id, created, status, lastUsedDate, lastUsedRegion, lastUsedService string

		if v.AccessKeyId != nil {
			id = *v.AccessKeyId

			out, err := i.iamClient.GetAccessKeyLastUsed(
				context.TODO(),
				&iam.GetAccessKeyLastUsedInput{
					AccessKeyId: aws.String(*v.AccessKeyId),
				},
			)
			if err == nil && out.AccessKeyLastUsed != nil {
				if out.AccessKeyLastUsed.LastUsedDate != nil {
					lastUsedDate = out.AccessKeyLastUsed.LastUsedDate.Format(utils.DefaultTimeFormat)
				}
				if out.AccessKeyLastUsed.Region != nil {
					lastUsedRegion = *out.AccessKeyLastUsed.Region
				}
				if out.AccessKeyLastUsed.ServiceName != nil {
					lastUsedService = *out.AccessKeyLastUsed.ServiceName
				}
			}
		}
		if v.CreateDate != nil {
			created = v.CreateDate.Format(utils.DefaultTimeFormat)
		}
		status = string(v.Status)

		data = append(data, []string{
			id,
			created,
			status,
			lastUsedDate,
			lastUsedRegion,
			lastUsedService,
		})
	}
	i.SetData(data)
}
