package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type IAMAccessKeys struct {
	*ui.Table
	repo     *repo.IAM
	userName string
	app      *Application
}

func NewIAMAccessKeys(repo *repo.IAM, userName string, app *Application) *IAMAccessKeys {
	i := &IAMAccessKeys{
		Table: ui.NewTable([]string{
			"ID",
			"CREATED",
			"STATUS",
			"LAST USED DATE",
			"LAST USED REGION",
			"LAST USED SERVICE",
		}, 1, 0),
		repo:     repo,
		userName: userName,
		app:      app,
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
	model, err := i.repo.ListAccessKeys(i.userName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
