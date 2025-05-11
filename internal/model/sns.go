package model

import (
	snsTypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type (
	SNSTopic struct {
		Arn        string
		Attributes map[string]string
	}
	SNSSubscription struct {
		snsTypes.Subscription
		Attributes map[string]string
	}
)
