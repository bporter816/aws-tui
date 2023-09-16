package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type CFTags struct {
	*ui.Table
	repo *repo.Cloudfront
	id   string
	app  *Application
}

func NewCFTags(repo *repo.Cloudfront, id string, app *Application) *CFTags {
	c := &CFTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		id:   id,
		app:  app,
	}
	return c
}

func (c CFTags) GetService() string {
	return "Cloudfront"
}
func (c CFTags) GetLabels() []string {
	// TODO generalize for other resources
	// extract id from arn
	arn, err := arn.Parse(c.id)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (c CFTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFTags) Render() {
	model, err := c.repo.ListTags(c.id)
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
