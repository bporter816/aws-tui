package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type EC2SecurityGroupTags struct {
	*ui.Table
	repo *repo.EC2
	sgId string
	app  *Application
}

func NewEC2SecurityGroupTags(repo *repo.EC2, sgId string, app *Application) *EC2SecurityGroupTags {
	e := &EC2SecurityGroupTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		sgId: sgId,
		app:  app,
	}
	return e
}

func (e EC2SecurityGroupTags) GetService() string {
	return "EC2"
}

func (e EC2SecurityGroupTags) GetLabels() []string {
	return []string{e.sgId, "Tags"}
}

func (e EC2SecurityGroupTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2SecurityGroupTags) Render() {
	model, err := e.repo.ListSecurityGroupTags(e.sgId)
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
	e.SetData(data)
}
