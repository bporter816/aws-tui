package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewIAMAccessKeys(repo *repo.IAM, userName string, app *Application) *TableView[model.IAMAccessKey] {
	return NewTableView(
		"IAM",
		[]string{userName, "Access Keys"},
		[]string{"ID", "CREATED", "STATUS", "LAST USED DATE", "LAST USED REGION", "LAST USED SERVICE"},
		func() ([]model.IAMAccessKey, error) {
			return repo.ListAccessKeys(userName)
		},
		func(v model.IAMAccessKey) []string {
			var id, created, status, lastUsedDate, lastUsedRegion, lastUsedService string

			if v.AccessKeyId != nil {
				id = *v.AccessKeyId

				if v.LastUsed.LastUsedDate != nil {
					lastUsedDate = v.LastUsed.LastUsedDate.Format(utils.DefaultTimeFormat)
				}
				if v.LastUsed.Region != nil {
					lastUsedRegion = *v.LastUsed.Region
				}
				if v.LastUsed.ServiceName != nil {
					lastUsedService = *v.LastUsed.ServiceName
				}
			}
			if v.CreateDate != nil {
				created = v.CreateDate.Format(utils.DefaultTimeFormat)
			}
			status = string(v.Status)

			return []string{
				id,
				created,
				status,
				lastUsedDate,
				lastUsedRegion,
				lastUsedService,
			}
		},
		nil,
	)
}
