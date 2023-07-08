package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

type DynamoDBTables struct {
	*Table
	ddbClient *ddb.Client
	app       *Application
	arns      []string
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
	return d
}

func (d DynamoDBTables) GetService() string {
	return "DynamoDB"
}

func (d DynamoDBTables) GetLabels() []string {
	return []string{"Tables"}
}

func (d DynamoDBTables) indexesHandler() {
	// TODO check if any indexes exist
	tableName, err := d.GetColSelection("NAME")
	if err != nil {
		return
	}
	indexesView := NewDynamoDBTableIndexes(d.ddbClient, tableName)
	d.app.AddAndSwitch(indexesView)
}

func (d DynamoDBTables) tagsHandler() {
	row, err := d.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewDynamoDBTags(d.ddbClient, d.arns[row-1], d.app)
	d.app.AddAndSwitch(tagsView)
}

func (d DynamoDBTables) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
			Description: "Indexes",
			Action:      d.indexesHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      d.tagsHandler,
		},
	}
}

func (d *DynamoDBTables) Render() {
	pg := ddb.NewListTablesPaginator(
		d.ddbClient,
		&ddb.ListTablesInput{},
	)
	var tableNames []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		tableNames = append(tableNames, out.TableNames...)
	}

	caser := cases.Title(language.English)
	var data [][]string
	d.arns = make([]string, len(tableNames))
	for i, v := range tableNames {
		out, err := d.ddbClient.DescribeTable(
			context.TODO(),
			&ddb.DescribeTableInput{
				TableName: aws.String(v),
			},
		)
		if err != nil {
			panic(err)
		}
		d.arns[i] = *out.Table.TableArn
		partitionKey, sortKey := getPartitionAndSortKeys(out.Table.KeySchema)
		if partitionKeyType, ok := getAttributeType(partitionKey, out.Table.AttributeDefinitions); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := getAttributeType(sortKey, out.Table.AttributeDefinitions); ok {
			sortKey = fmt.Sprintf("%v (%v)", sortKey, sortKeyType)
		}
		var billingMode ddbTypes.BillingMode
		if out.Table.BillingModeSummary != nil {
			billingMode = out.Table.BillingModeSummary.BillingMode
		} else {
			// tables that have never had on-demand capacity set appear to not return this part of the response at all
			billingMode = ddbTypes.BillingModeProvisioned
		}
		var readCap, writeCap string
		var itemCount, tableSize int64
		if billingMode == ddbTypes.BillingModePayPerRequest {
			readCap, writeCap = "-", "-"
		} else if out.Table.ProvisionedThroughput != nil {
			if out.Table.ProvisionedThroughput.ReadCapacityUnits != nil {
				readCap = strconv.FormatInt(*out.Table.ProvisionedThroughput.ReadCapacityUnits, 10)
			}
			if out.Table.ProvisionedThroughput.WriteCapacityUnits != nil {
				writeCap = strconv.FormatInt(*out.Table.ProvisionedThroughput.WriteCapacityUnits, 10)
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
			caser.String(string(out.Table.TableStatus)),
			partitionKey,
			sortKey,
			strconv.Itoa(len(out.Table.GlobalSecondaryIndexes) + len(out.Table.LocalSecondaryIndexes)),
			caser.String(strings.ReplaceAll(string(billingMode), "_", " ")),
			readCap,
			writeCap,
			strconv.FormatInt(itemCount, 10),
			strconv.FormatInt(tableSize, 10),
		})
	}
	d.SetData(data)
}
