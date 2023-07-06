package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"strings"
)

type ElasticacheTags struct {
	*Table
	ecClient     *ec.Client
	resourceName string
	app          *Application
}

func NewElasticacheTags(ecClient *ec.Client, resourceName string, app *Application) *ElasticacheTags {
	e := &ElasticacheTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ecClient:     ecClient,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ElasticacheTags) GetName() string {
	// TODO generalize for other resources
	// extract id from arn
	parts := strings.Split(e.resourceName, ":")
	id := parts[len(parts)-1]
	return fmt.Sprintf("Elasticache | Clusters | %v | Tags", id)
}

func (e ElasticacheTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheTags) Render() {
	out, err := e.ecClient.ListTagsForResource(
		context.TODO(),
		&ec.ListTagsForResourceInput{
			ResourceName: aws.String(e.resourceName),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range out.TagList {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	e.SetData(data)
}
