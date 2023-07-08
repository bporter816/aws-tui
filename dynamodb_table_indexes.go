package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
	"strings"
)

type DynamoDBTableIndexes struct {
	*Table
	ddbClient *ddb.Client
	tableName string
}

func NewDynamoDBTableIndexes(ddbClient *ddb.Client, tableName string) *DynamoDBTableIndexes {
	d := &DynamoDBTableIndexes{
		Table: NewTable([]string{
			"NAME",
			"TYPE",
			"STATUS",
			"PARTITION KEY",
			"SORT KEY",
			"ITEMS",
			"SIZE",
			"PROJECTED ATTRIBUTES",
		}, 1, 0),
		ddbClient: ddbClient,
		tableName: tableName,
	}
	return d
}

func (d DynamoDBTableIndexes) GetService() string {
	return "DynamoDB"
}

func (d DynamoDBTableIndexes) GetLabels() []string {
	return []string{d.tableName, "Indexes"}
}

func (d DynamoDBTableIndexes) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (d DynamoDBTableIndexes) Render() {
	out, err := d.ddbClient.DescribeTable(
		context.TODO(),
		&ddb.DescribeTableInput{
			TableName: aws.String(d.tableName),
		},
	)
	if err != nil {
		panic(err)
	}
	var data [][]string
	for _, v := range out.Table.GlobalSecondaryIndexes {
		partitionKey, sortKey := getPartitionAndSortKeys(v.KeySchema)
		if partitionKeyType, ok := getAttributeType(partitionKey, out.Table.AttributeDefinitions); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := getAttributeType(sortKey, out.Table.AttributeDefinitions); ok {
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
	for _, v := range out.Table.LocalSecondaryIndexes {
		partitionKey, sortKey := getPartitionAndSortKeys(v.KeySchema)
		if partitionKeyType, ok := getAttributeType(partitionKey, out.Table.AttributeDefinitions); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := getAttributeType(sortKey, out.Table.AttributeDefinitions); ok {
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
