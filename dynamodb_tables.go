package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

type DynamoDBTables struct {
	*Table
	ddbClient *ddb.Client
	app       *Application
}

func NewDynamoDBTables(ddbClient *ddb.Client, app *Application) *DynamoDBTables {
	d := &DynamoDBTables{
		Table: NewTable([]string{
			"NAME",
			"STATUS",
			"PARTITION KEY",
			"SORT KEY",
			"INDEXES",
			"BILLING",
			"READ CAP",
			"WRITE CAP",
			"ITEMS",
			"SIZE",
		}, 1, 0),
		ddbClient: ddbClient,
		app:       app,
	}
	d.Render() // TODO fix
	return d
}

func (d DynamoDBTables) GetName() string {
	return "DynamoDB"
}

func (d DynamoDBTables) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (d DynamoDBTables) Render() {
	tablesPaginator := ddb.NewListTablesPaginator(d.ddbClient, &ddb.ListTablesInput{})
	var tableNames []string
	for tablesPaginator.HasMorePages() {
		out, err := tablesPaginator.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		tableNames = append(tableNames, out.TableNames...)
	}

	var data [][]string
	for _, v := range tableNames {
		out, err := d.ddbClient.DescribeTable(
			context.TODO(),
			&ddb.DescribeTableInput{
				TableName: aws.String(v),
			},
		)
		if err != nil {
			panic(err)
		}
		var partitionKey, sortKey string
		for _, ks := range out.Table.KeySchema {
			if ks.KeyType == ddbTypes.KeyTypeHash {
				partitionKey = *ks.AttributeName
			} else {
				// TODO explicitly check if this is a range key?
				sortKey = *ks.AttributeName
			}
		}
		partitionKeyType, err := getAttributeType(partitionKey, out.Table.AttributeDefinitions)
		if err != nil {
			panic(err)
		}
		var billingMode string
		if out.Table.BillingModeSummary != nil {
			billingMode = string(out.Table.BillingModeSummary.BillingMode)
		}
		var readCap, writeCap, itemCount, tableSize int64
		if out.Table.ProvisionedThroughput != nil {
			if out.Table.ProvisionedThroughput.ReadCapacityUnits != nil {
				readCap = *out.Table.ProvisionedThroughput.ReadCapacityUnits
			}
			if out.Table.ProvisionedThroughput.WriteCapacityUnits != nil {
				writeCap = *out.Table.ProvisionedThroughput.WriteCapacityUnits
			}
		}
		if out.Table.ItemCount != nil {
			itemCount = *out.Table.ItemCount
		}
		if out.Table.TableSizeBytes != nil {
			tableSize = *out.Table.TableSizeBytes
		}
		data = append(data, []string{
			v,
			string(out.Table.TableStatus),
			fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType),
			sortKey,
			strconv.Itoa(len(out.Table.GlobalSecondaryIndexes) + len(out.Table.LocalSecondaryIndexes)),
			billingMode,
			strconv.FormatInt(readCap, 10),
			strconv.FormatInt(writeCap, 10),
			strconv.FormatInt(itemCount, 10),
			strconv.FormatInt(tableSize, 10),
		})
	}
	d.SetData(data)
}

func getAttributeType(attribute string, defs []ddbTypes.AttributeDefinition) (string, error) {
	for _, v := range defs {
		if attribute == *v.AttributeName {
			return string(v.AttributeType), nil
		}
	}
	return "", errors.New("no matching attribute")
}
