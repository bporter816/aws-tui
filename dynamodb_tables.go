package main

import (
	"fmt"
	"strconv"
	"strings"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type DynamoDBTables struct {
	*ui.Table
	view.DynamoDB
	repo  *repo.DynamoDB
	app   *Application
	model []model.DynamoDBTable
}

func NewDynamoDBTables(repo *repo.DynamoDB, app *Application) *DynamoDBTables {
	d := &DynamoDBTables{
		Table: ui.NewTable([]string{
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
		repo: repo,
		app:  app,
	}
	return d
}

func (d DynamoDBTables) GetLabels() []string {
	return []string{"Tables"}
}

func (d DynamoDBTables) indexesHandler() {
	// TODO check if any indexes exist
	row, err := d.GetRowSelection()
	if err != nil {
		return
	}
	tableName, err := d.GetColSelection("NAME")
	if err != nil {
		return
	}
	indexesView := NewDynamoDBTableIndexes(d.repo, tableName, d.model[row-1].AttributeDefinitions, d.app)
	d.app.AddAndSwitch(indexesView)
}

func (d DynamoDBTables) tagsHandler() {
	row, err := d.GetRowSelection()
	if err != nil {
		return
	}
	if d.model[row-1].TableArn == nil {
		return
	}
	tagsView := NewTags(d.repo, d.GetService(), *d.model[row-1].TableArn, d.app)
	d.app.AddAndSwitch(tagsView)
}

func (d DynamoDBTables) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
			Description: "Indexes",
			Action:      d.indexesHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      d.tagsHandler,
		},
	}
}

func (d *DynamoDBTables) Render() {
	model, err := d.repo.ListTables()
	if err != nil {
		panic(err)
	}
	d.model = model

	var data [][]string
	for _, v := range model {
		var name string
		if v.TableName != nil {
			name = *v.TableName
		}
		partitionKey, sortKey := utils.GetDynamoDBPartitionAndSortKeys(v.KeySchema)
		if partitionKeyType, ok := utils.GetDynamoDBAttributeType(partitionKey, v.AttributeDefinitions); ok {
			partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
		}
		if sortKeyType, ok := utils.GetDynamoDBAttributeType(sortKey, v.AttributeDefinitions); ok {
			sortKey = fmt.Sprintf("%v (%v)", sortKey, sortKeyType)
		}
		var billingMode ddbTypes.BillingMode
		if v.BillingModeSummary != nil {
			billingMode = v.BillingModeSummary.BillingMode
		} else {
			// tables that have never had on-demand capacity set appear to not return this part of the response at all
			billingMode = ddbTypes.BillingModeProvisioned
		}
		var readCap, writeCap string
		var itemCount, tableSize int64
		if billingMode == ddbTypes.BillingModePayPerRequest {
			readCap, writeCap = "-", "-"
		} else if v.ProvisionedThroughput != nil {
			if v.ProvisionedThroughput.ReadCapacityUnits != nil {
				readCap = strconv.FormatInt(*v.ProvisionedThroughput.ReadCapacityUnits, 10)
			}
			if v.ProvisionedThroughput.WriteCapacityUnits != nil {
				writeCap = strconv.FormatInt(*v.ProvisionedThroughput.WriteCapacityUnits, 10)
			}
		}
		if v.ItemCount != nil {
			itemCount = *v.ItemCount
		}
		if v.TableSizeBytes != nil {
			tableSize = *v.TableSizeBytes
		}
		data = append(data, []string{
			name,
			utils.TitleCase(string(v.TableStatus)),
			partitionKey,
			sortKey,
			strconv.Itoa(len(v.GlobalSecondaryIndexes) + len(v.LocalSecondaryIndexes)),
			utils.TitleCase(strings.ReplaceAll(string(billingMode), "_", " ")),
			readCap,
			writeCap,
			strconv.FormatInt(itemCount, 10),
			utils.FormatSize(tableSize, 1),
		})
	}
	d.SetData(data)
}
