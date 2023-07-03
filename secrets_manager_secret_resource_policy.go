package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerSecretResourcePolicy struct {
	*Text
	smClient *sm.Client
	secretId string
}

func NewSecretsManagerSecretResourcePolicy(smClient *sm.Client, secretId string) *SecretsManagerSecretResourcePolicy {
	s := &SecretsManagerSecretResourcePolicy{
		Text:     NewText(true, "json"),
		smClient: smClient,
		secretId: secretId,
	}
	return s
}

func (s SecretsManagerSecretResourcePolicy) GetName() string {
	return fmt.Sprintf("Secrets Manager | %v | Resource Policy", s.secretId)
}

func (s SecretsManagerSecretResourcePolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SecretsManagerSecretResourcePolicy) Render() {
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
