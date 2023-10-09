package model

import (
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type (
	KMSKey struct {
		kmsTypes.KeyMetadata
		Aliases []string
	}
	KMSGrant          kmsTypes.GrantListEntry
	KMSCustomKeyStore kmsTypes.CustomKeyStoresListEntry
)
