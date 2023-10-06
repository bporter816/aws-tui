package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type CloudwatchTags struct {
	*ui.Table
	repo *repo.Cloudwatch
	arn  string
	name string
	app  *Application
}

func NewCloudwatchTags(repo *repo.Cloudwatch, arn string, name string, app *Application) *CloudwatchTags {
	c := &CloudwatchTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		arn:  arn,
		name: name,
		app:  app,
	}
	return c
}

func (c CloudwatchTags) GetService() string {
	return "Cloudwatch"
}

func (c CloudwatchTags) GetLabels() []string {
	return []string{c.name, "Tags"}
}

func (c CloudwatchTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CloudwatchTags) Render() {
	model, err := c.repo.ListTags(c.arn)
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
	c.SetData(data)
}
