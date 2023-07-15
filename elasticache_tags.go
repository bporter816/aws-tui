package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ec "github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/bporter816/aws-tui/ui"
)

type ElasticacheTags struct {
	*ui.Table
	ecClient     *ec.Client
	resourceType ElasticacheResourceType
	resourceArn  string
	resourceName string
	app          *Application
}

type ElasticacheResourceType string

const (
	ElasticacheResourceTypeCluster      ElasticacheResourceType = "Clusters"
	ElasticacheResourceTypeReservedNode ElasticacheResourceType = "Reserved Nodes"
	ElasticacheResourceTypeSnapshot     ElasticacheResourceType = "Snapshots"
)

func NewElasticacheTags(ecClient *ec.Client, resourceType ElasticacheResourceType, resourceArn string, resourceName string, app *Application) *ElasticacheTags {
	e := &ElasticacheTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		ecClient:     ecClient,
		resourceType: resourceType,
		resourceArn:  resourceArn,
		resourceName: resourceName,
		app:          app,
	}
	return e
}

func (e ElasticacheTags) GetService() string {
	return "Elasticache"
}

func (e ElasticacheTags) GetLabels() []string {
	return []string{e.resourceName, "Tags"}
}

func (e ElasticacheTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ElasticacheTags) Render() {
	out, err := e.ecClient.ListTagsForResource(
		context.TODO(),
		&ec.ListTagsForResourceInput{
			ResourceName: aws.String(e.resourceArn),
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
