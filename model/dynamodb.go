package model

import (
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type (
	DynamoDBTable ddbTypes.TableDescription
	DynamoDBIndexes struct {
		Global []ddbTypes.GlobalSecondaryIndexDescription
		Local []ddbTypes.LocalSecondaryIndexDescription
	}
)
