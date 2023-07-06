package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SMSecretResourcePolicy struct {
	*Text
	smClient *sm.Client
	secretId string
}

func NewSMSecretResourcePolicy(smClient *sm.Client, secretId string) *SMSecretResourcePolicy {
	s := &SMSecretResourcePolicy{
		Text:     NewText(true, "json"),
		smClient: smClient,
		secretId: secretId,
	}
	return s
}

func (s SMSecretResourcePolicy) GetName() string {
	return fmt.Sprintf("Secrets Manager | %v | Resource Policy", s.secretId)
}

func (s SMSecretResourcePolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretResourcePolicy) Render() {
	out, err := s.smClient.GetResourcePolicy(
		context.TODO(),
		&sm.GetResourcePolicyInput{
			SecretId: aws.String(s.secretId),
		},
	)
	if err != nil {
		panic(err)
	}
	var policy string
	if out.ResourcePolicy != nil {
		policy = *out.ResourcePolicy
	}
	s.SetText(policy)
}
