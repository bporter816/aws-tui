package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type CloudWatchTags struct {
	*ui.Table
	repo *repo.CloudWatch
	arn  string
	name string
	app  *Application
}

func NewCloudWatchTags(repo *repo.CloudWatch, arn string, name string, app *Application) *CloudWatchTags {
	c := &CloudWatchTags{
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

func (c CloudWatchTags) GetService() string {
	return "CloudWatch"
}

func (c CloudWatchTags) GetLabels() []string {
	return []string{c.name, "Tags"}
}

func (c CloudWatchTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CloudWatchTags) Render() {
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
