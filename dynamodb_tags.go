package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/bporter816/aws-tui/ui"
	"strings"
)

type DynamoDBTags struct {
	*ui.Table
	ddbClient *ddb.Client
	id        string
	app       *Application
}

func NewDynamoDBTags(ddbClient *ddb.Client, id string, app *Application) *DynamoDBTags {
	d := &DynamoDBTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ddbClient: ddbClient,
		id:        id,
		app:       app,
	}
	return d
}

func (d DynamoDBTags) GetService() string {
	return "Cloudfront"
}

func (d DynamoDBTags) GetLabels() []string {
	// TODO generalize for other resources
	// extract id from arn
	parts := strings.Split(d.id, "/")
	id := parts[len(parts)-1]
	return []string{id, "Tags"}
}

func (d DynamoDBTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (d DynamoDBTags) Render() {
	out, err := d.ddbClient.ListTagsOfResource(
		context.TODO(),
		&ddb.ListTagsOfResourceInput{
			ResourceArn: aws.String(d.id),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.Tags != nil {
		for _, v := range out.Tags {
			data = append(data, []string{
				*v.Key,
				*v.Value,
			})
		}
	}
	d.SetData(data)
}
