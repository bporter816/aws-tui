package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"strings"
)

type CFTags struct {
	*Table
	cfClient *cf.Client
	id       string
	app      *Application
}

func NewCFTags(cfClient *cf.Client, id string, app *Application) *CFTags {
	c := &CFTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		cfClient: cfClient,
		id:       id,
		app:      app,
	}
	return c
}

func (c CFTags) GetService() string {
	return "Cloudfront"
}
func (c CFTags) GetLabels() []string {
	// TODO generalize for other resources
	// extract id from arn
	parts := strings.Split(c.id, "/")
	id := parts[len(parts)-1]
	return []string{id, "Tags"}
}

func (c CFTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (c CFTags) Render() {
	out, err := c.cfClient.ListTagsForResource(
		context.TODO(),
		&cf.ListTagsForResourceInput{
			Resource: aws.String(c.id),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	if out.Tags != nil {
		for _, v := range out.Tags.Items {
			data = append(data, []string{
				*v.Key,
				*v.Value,
			})
		}
	}
	c.SetData(data)
}
