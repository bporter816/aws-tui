package internal

import (
	"fmt"
	"strconv"
	"strings"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

type dynamoDBIndex struct {
	name       string
	indexType  string
	status     string
	keySchema  []ddbTypes.KeySchemaElement
	itemCount  int64
	sizeBytes  int64
	projection *ddbTypes.Projection
}

func NewDynamoDBTableIndexes(repo *repo.DynamoDB, tableName string, attributes []ddbTypes.AttributeDefinition, app *Application) *TableView[dynamoDBIndex] {
	return NewTableView(
		"DynamoDB",
		[]string{tableName, "Indexes"},
		[]string{"NAME", "TYPE", "STATUS", "PARTITION KEY", "SORT KEY", "ITEMS", "SIZE", "PROJECTED ATTRIBUTES"},
		func() ([]dynamoDBIndex, error) {
			model, err := repo.ListIndexes(tableName)
			if err != nil {
				return nil, err
			}
			var indexes []dynamoDBIndex
			for _, v := range model.Global {
				var itemCount, sizeBytes int64
				if v.ItemCount != nil {
					itemCount = *v.ItemCount
				}
				if v.IndexSizeBytes != nil {
					sizeBytes = *v.IndexSizeBytes
				}
				indexes = append(indexes, dynamoDBIndex{
					name:       utils.DerefString(v.IndexName, ""),
					indexType:  "Global",
					status:     string(v.IndexStatus),
					keySchema:  v.KeySchema,
					itemCount:  itemCount,
					sizeBytes:  sizeBytes,
					projection: v.Projection,
				})
			}
			for _, v := range model.Local {
				var itemCount, sizeBytes int64
				if v.ItemCount != nil {
					itemCount = *v.ItemCount
				}
				if v.IndexSizeBytes != nil {
					sizeBytes = *v.IndexSizeBytes
				}
				indexes = append(indexes, dynamoDBIndex{
					name:       utils.DerefString(v.IndexName, ""),
					indexType:  "Local",
					status:     "-",
					keySchema:  v.KeySchema,
					itemCount:  itemCount,
					sizeBytes:  sizeBytes,
					projection: v.Projection,
				})
			}
			return indexes, nil
		},
		func(v dynamoDBIndex) []string {
			partitionKey, sortKey := utils.GetDynamoDBPartitionAndSortKeys(v.keySchema)
			if partitionKeyType, ok := utils.GetDynamoDBAttributeType(partitionKey, attributes); ok {
				partitionKey = fmt.Sprintf("%v (%v)", partitionKey, partitionKeyType)
			}
			if sortKeyType, ok := utils.GetDynamoDBAttributeType(sortKey, attributes); ok {
				sortKey = fmt.Sprintf("%v (%v)", sortKey, sortKeyType)
			}
			var projection string
			if v.projection != nil {
				if v.projection.ProjectionType == ddbTypes.ProjectionTypeInclude {
					projection = strings.Join(v.projection.NonKeyAttributes, ", ")
				} else {
					projection = string(v.projection.ProjectionType)
				}
			}
			return []string{
				v.name,
				v.indexType,
				v.status,
				partitionKey,
				sortKey,
				strconv.FormatInt(v.itemCount, 10),
				strconv.FormatInt(v.sizeBytes, 10),
				projection,
			}
		},
		nil,
	)
}
