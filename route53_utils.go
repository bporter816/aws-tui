package main

import (
	"errors"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func getTag(tags []r53Types.Tag, key string) (string, error) {
	for _, v := range tags {
		if *v.Key == key {
			return *v.Value, nil
		}
	}
	return "", errors.New("tag key not found")
}
