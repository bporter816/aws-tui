package model

import (
	smTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

type SecretsManagerSecret smTypes.SecretListEntry
