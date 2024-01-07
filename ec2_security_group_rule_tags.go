package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type EC2SecurityGroupRuleTags struct {
	*ui.Table
	view.EC2
	repo   *repo.EC2
	ruleId string
	app    *Application
}

func NewEC2SecurityGroupRuleTags(repo *repo.EC2, ruleId string, app *Application) *EC2SecurityGroupRuleTags {
	e := &EC2SecurityGroupRuleTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:   repo,
		ruleId: ruleId,
		app:    app,
	}
	return e
}

func (e EC2SecurityGroupRuleTags) GetLabels() []string {
	return []string{e.ruleId, "Tags"}
}

func (e EC2SecurityGroupRuleTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2SecurityGroupRuleTags) Render() {
	model, err := e.repo.ListSecurityGroupRuleTags(e.ruleId)
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
