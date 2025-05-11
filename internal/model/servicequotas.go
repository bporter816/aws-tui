package model

import (
	sqTypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
)

type (
	ServiceQuotasService sqTypes.ServiceInfo

	// There are different APIs to return default and applied quotas, so this consolidates them
	ServiceQuotasQuota struct {
		sqTypes.ServiceQuota
		DefaultValue *float64
		AppliedValue *float64
	}
)
