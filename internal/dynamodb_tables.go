package internal

import (
	"fmt"
	"strconv"
	"strings"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewDynamoDBTables(repo *repo.DynamoDB, app *Application) *TableView[model.DynamoDBTable] {
	return NewTableView(
		"DynamoDB",
		[]string{"Tables"},
		[]string{"NAME", "STATUS", "PARTITION KEY", "SORT KEY", "INDEXES", "BILLING", "READ CAP", "WRITE CAP", "ITEMS", "SIZE"},
		repo.ListTables,
		func(v model.DynamoDBTable) []string {
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
			return []string{
				utils.DerefString(v.TableName, ""),
				utils.TitleCase(string(v.TableStatus)),
				partitionKey,
				sortKey,
				strconv.Itoa(len(v.GlobalSecondaryIndexes) + len(v.LocalSecondaryIndexes)),
				utils.TitleCase(strings.ReplaceAll(string(billingMode), "_", " ")),
				readCap,
				writeCap,
				strconv.FormatInt(itemCount, 10),
				utils.FormatSize(tableSize, 1),
			}
		},
		func(v *TableView[model.DynamoDBTable]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
					Description: "Indexes",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						tableName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewDynamoDBTableIndexes(repo, tableName, item.AttributeDefinitions, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TableArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.TableArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("DynamoDB.Tables", func(r Repos, app *Application) Component {
		return NewDynamoDBTables(r.DynamoDB, app)
	})
}
