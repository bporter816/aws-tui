package repo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	ddb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bporter816/aws-tui/model"
)

type DynamoDB struct {
	ddbClient *ddb.Client
}

func NewDynamoDB(ddbClient *ddb.Client) *DynamoDB {
	return &DynamoDB{
		ddbClient: ddbClient,
	}
}

func (d DynamoDB) listTableNames() ([]string, error) {
	pg := ddb.NewListTablesPaginator(
		d.ddbClient,
		&ddb.ListTablesInput{},
	)
	var tableNames []string
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			return []string{}, err
		}
		tableNames = append(tableNames, out.TableNames...)
	}
	return tableNames, nil
}

func (d DynamoDB) describeTable(tableName string) (ddbTypes.TableDescription, error) {
	out, err := d.ddbClient.DescribeTable(
		context.TODO(),
		&ddb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
	)
	if err != nil || out.Table == nil {
		return ddbTypes.TableDescription{}, err
	}
	return *out.Table, nil
}

func (d DynamoDB) ListTables() ([]model.DynamoDBTable, error) {
	tableNames, err := d.listTableNames()
	if err != nil {
		return []model.DynamoDBTable{}, err
	}

	var tables []model.DynamoDBTable
	for _, v := range tableNames {
		table, err := d.describeTable(v)
		if err != nil {
			// TODO handle error more cleanly
			tables = append(tables, model.DynamoDBTable{TableName: aws.String(v)})
		} else {
			tables = append(tables, model.DynamoDBTable(table))
		}
	}
	return tables, nil
}

func (d DynamoDB) ListIndexes(tableName string) (model.DynamoDBIndexes, error) {
	table, err := d.describeTable(tableName)
	if err != nil {
		return model.DynamoDBIndexes{}, err
	}
	return model.DynamoDBIndexes{Global: table.GlobalSecondaryIndexes, Local: table.LocalSecondaryIndexes}, nil
}

func (d DynamoDB) ListTags(resourceId string) (model.Tags, error) {
	out, err := d.ddbClient.ListTagsOfResource(
		context.TODO(),
		&ddb.ListTagsOfResourceInput{
			ResourceArn: aws.String(resourceId),
		},
	)
	if err != nil {
		return model.Tags{}, err
	}
	var tags model.Tags
	for _, v := range out.Tags {
		tags = append(tags, model.Tag{Key: *v.Key, Value: *v.Value})
	}
	return tags, nil
}
