package model

import (
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type (
	ELBLoadBalancer elbTypes.LoadBalancer
	ELBListener     struct {
		elbTypes.Listener
		Rules int
	}
	ELBListenerRule          elbTypes.Rule
	ELBTargetGroup           elbTypes.TargetGroup
	ELBTrustStore            elbTypes.TrustStore
	ELBTrustStoreAssociation elbTypes.TrustStoreAssociation
)
