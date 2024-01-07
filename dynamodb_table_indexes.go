package main

import (
	"fmt"
	"strconv"
	"strings"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type DynamoDBTableIndexes struct {
	*ui.Table
	view.DynamoDB
	repo       *repo.DynamoDB
	tableName  string
	attributes []ddbTypes.AttributeDefinition
	app        *Application
}

func NewDynamoDBTableIndexes(repo *repo.DynamoDB, tableName string, attributes []ddbTypes.AttributeDefinition, app *Application) *DynamoDBTableIndexes {
	d := &DynamoDBTableIndexes{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"STATUS",
			"PARTITION KEY",
			"SORT KEY",
			"ITEMS",
			"SIZE",
			"PROJECTED ATTRIBUTES",
		}, 1, 0),
		repo:       repo,
		tableName:  tableName,
		attributes: attributes,
		app:        app,
	}
	return d
}

func (d DynamoDBTableIndexes) GetLabels() []string {
	return []string{d.tableName, "Indexes"}
}

func (d DynamoDBTableIndexes) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (d DynamoDBTableIndexes) Render() {
	model, err := d.repo.ListIndexes(d.tableName)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model.Global {
		partitionKey, sortKey := getPartitionAndSortKeys(v.KeySchema)
		if partitionKeyType, ok := getAttributeType(partitionKey, d.attributes); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := getAttributeType(sortKey, d.attributes); ok {
			sortKey = fmt.Sprintf("%v (%v)", sortKey, sortKeyType)
		}
		var projection string
		if v.Projection != nil {
			if v.Projection.ProjectionType == ddbTypes.ProjectionTypeInclude {
				projection = strings.Join(v.Projection.NonKeyAttributes, ", ")
			} else {
				projection = string(v.Projection.ProjectionType)
			}
		}
		data = append(data, []string{
			*v.IndexName,
			"Global",
			string(v.IndexStatus),
			partitionKey,
			sortKey,
			strconv.FormatInt(*v.ItemCount, 10),
			strconv.FormatInt(*v.IndexSizeBytes, 10),
			projection,
		})
	}
	for _, v := range model.Local {
		partitionKey, sortKey := getPartitionAndSortKeys(v.KeySchema)
		if partitionKeyType, ok := getAttributeType(partitionKey, d.attributes); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := getAttributeType(sortKey, d.attributes); ok {
			sortKey = fmt.Sprintf("%v (%v)", sortKey, sortKeyType)
		}
		var projection string
		if v.Projection != nil {
			if v.Projection.ProjectionType == ddbTypes.ProjectionTypeInclude {
				projection = strings.Join(v.Projection.NonKeyAttributes, ", ")
			} else {
				projection = string(v.Projection.ProjectionType)
			}
		}
		data = append(data, []string{
			*v.IndexName,
			"Local",
			"-",
			partitionKey,
			sortKey,
			strconv.FormatInt(*v.ItemCount, 10),
			strconv.FormatInt(*v.IndexSizeBytes, 10),
			projection,
		})
	}
	d.SetData(data)
}
