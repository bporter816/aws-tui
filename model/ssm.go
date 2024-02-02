package model

import (
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type (
	SSMParameter ssmTypes.ParameterMetadata
)
