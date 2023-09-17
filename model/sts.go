package model

import (
	sts "github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSCallerIdentity sts.GetCallerIdentityOutput
