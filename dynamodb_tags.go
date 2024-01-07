package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type DynamoDBTags struct {
	*ui.Table
	view.DynamoDB
	repo *repo.DynamoDB
	id   string
	app  *Application
}

func NewDynamoDBTags(repo *repo.DynamoDB, id string, app *Application) *DynamoDBTags {
	d := &DynamoDBTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo: repo,
		id:   id,
		app:  app,
	}
	return d
}

func (d DynamoDBTags) GetLabels() []string {
	// extract id from arn
	arn, err := arn.Parse(d.id)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (d DynamoDBTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (d DynamoDBTags) Render() {
	model, err := d.repo.ListTags(d.id)
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
	d.SetData(data)
}
