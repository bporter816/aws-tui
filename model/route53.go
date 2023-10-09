package model

import (
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type (
	Route53HostedZone  r53Types.HostedZone
	Route53Record      r53Types.ResourceRecordSet
	Route53HealthCheck r53Types.HealthCheck
)
