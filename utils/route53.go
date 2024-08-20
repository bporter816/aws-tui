package utils

import (
	"fmt"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"strings"
)

// TODO add tests
func FormatRoute53HealthCheckDescription(v r53Types.HealthCheckConfig) string {
	if v.Type == r53Types.HealthCheckTypeCalculated {
		return fmt.Sprintf("Calculated threshold: %v out of %v", *v.HealthThreshold, len(v.ChildHealthChecks))
	} else if v.Type == r53Types.HealthCheckTypeCloudwatchMetric {
		if v.AlarmIdentifier != nil && v.AlarmIdentifier.Name != nil {
			return fmt.Sprintf("CloudWatch alarm: %v (%v)", *v.AlarmIdentifier.Name, v.AlarmIdentifier.Region)
		}
		return ""
	} else if v.Type == r53Types.HealthCheckTypeRecoveryControl {
		// TODO handle this case?
		return ""
	} else {
		var protocol, host, path string
		var port int32
		if v.Type == r53Types.HealthCheckTypeHttp || v.Type == r53Types.HealthCheckTypeHttpStrMatch {
			protocol = "http"
		} else if v.Type == r53Types.HealthCheckTypeHttps || v.Type == r53Types.HealthCheckTypeHttpsStrMatch {
			protocol = "https"
		} else if v.Type == r53Types.HealthCheckTypeTcp {
			protocol = "tcp"
		}
		if v.FullyQualifiedDomainName != nil {
			host = *v.FullyQualifiedDomainName
		} else {
			return ""
		}
		if v.Port != nil {
			port = *v.Port
		} else {
			return ""
		}
		if v.ResourcePath != nil {
			path = *v.ResourcePath
		}
		return fmt.Sprintf("%v://%v:%v%v", protocol, host, port, path)
	}
}

// TODO add tests
func FormatRoute53ResourceRecords(items []r53Types.ResourceRecord) string {
	values := make([]string, len(items))
	for i, v := range items {
		if v.Value != nil {
			values[i] = string(*v.Value)
		}
	}
	return TruncateStrings(values, 1)
}

// TODO consider removing, also see route53_records.go
// TODO add tests
func JoinRoute53ResourceRecords(items []r53Types.ResourceRecord, sep string) string {
	if len(items) == 0 {
		return ""
	}
	var ret string
	for _, v := range items {
		ret += sep
		ret += *v.Value
	}
	return strings.TrimPrefix(ret, sep)
}
