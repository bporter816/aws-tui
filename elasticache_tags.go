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
	// extract resource type and id from arn
	parts := strings.Split(e.resourceName, ":")
	resourceType, id := parts[len(parts)-2], parts[len(parts)-1]
	var resourceTypeStr string
	switch resourceType {
	case "cluster":
		resourceTypeStr = "Clusters"
	case "reserved-instance":
		resourceTypeStr = "Reserved Nodes"
	case "snapshot":
		resourceTypeStr = "Snapshots"
	default:
		resourceTypeStr = "<unknown>"

	}
	return fmt.Sprintf("Elasticache | %v | %v | Tags", resourceTypeStr, id)
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
