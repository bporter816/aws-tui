package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type IAMRoleTags struct {
	*ui.Table
	repo     *repo.IAM
	roleName string
	app      *Application
}

func NewIAMRoleTags(repo *repo.IAM, roleName string, app *Application) *IAMRoleTags {
	i := &IAMRoleTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:     repo,
		roleName: roleName,
		app:      app,
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
	model, err := i.repo.ListRoleTags(i.roleName)
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
