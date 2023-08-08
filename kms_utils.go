package main

import (
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"strings"
)

func joinGrantOperations(operations []kmsTypes.GrantOperation, sep string) string {
	if len(operations) == 0 {
		return ""
	}
	var ret string
	for _, v := range operations {
		ret += sep
		ret += string(v)
	}
	return strings.TrimPrefix(ret, sep)
}
