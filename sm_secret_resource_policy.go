package main

import (
	"context"
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

func (s SMSecretResourcePolicy) GetService() string {
	return "Secrets Manager"
}

func (s SMSecretResourcePolicy) GetLabels() []string {
	return []string{s.secretId, "Resource Policy"}
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
