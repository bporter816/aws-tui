package utils

import (
	"fmt"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"strings"
)

func JoinKMSGrantOperations(operations []kmsTypes.GrantOperation, sep string) string {
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

func FormatKMSAliases(a []string) string {
	if len(a) == 0 {
		return "-"
	}
	if len(a) == 1 {
		return a[0]
	}
	return fmt.Sprintf("%v + %v more", a[0], len(a)-1)
}
