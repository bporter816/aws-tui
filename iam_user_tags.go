package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type IAMUserTags struct {
	*ui.Table
	view.IAM
	repo     *repo.IAM
	userName string
	app      *Application
}

func NewIAMUserTags(repo *repo.IAM, userName string, app *Application) *IAMUserTags {
	i := &IAMUserTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:     repo,
		userName: userName,
		app:      app,
	}
	return i
}

func (i IAMUserTags) GetLabels() []string {
	return []string{i.userName, "Tags"}
}

func (i IAMUserTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (i IAMUserTags) Render() {
	model, err := i.repo.ListUserTags(i.userName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	i.SetData(data)
}
