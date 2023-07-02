package main

import (
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func getPartitionAndSortKeys(keySchema []ddbTypes.KeySchemaElement) (string, string) {
	var partitionKey, sortKey string
	for _, ks := range keySchema {
		if ks.KeyType == ddbTypes.KeyTypeHash {
			partitionKey = *ks.AttributeName
		} else if ks.KeyType == ddbTypes.KeyTypeRange {
			sortKey = *ks.AttributeName
		}
	}
	return partitionKey, sortKey
}

func getAttributeType(attribute string, defs []ddbTypes.AttributeDefinition) (string, bool) {
	for _, v := range defs {
		if attribute == *v.AttributeName {
			return string(v.AttributeType), true
		}
	}
	return "", false
}
