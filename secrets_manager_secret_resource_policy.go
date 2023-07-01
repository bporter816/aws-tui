package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rivo/tview"
)

type SecretsManagerSecretResourcePolicy struct {
	*tview.TextView
	smClient *sm.Client
	secretId string
}

func NewSecretsManagerSecretResourcePolicy(smClient *sm.Client, secretId string) *SecretsManagerSecretResourcePolicy {
	s := &SecretsManagerSecretResourcePolicy{
		TextView: tview.NewTextView().SetDynamicColors(true),
		smClient: smClient,
		secretId: secretId,
	}
	s.Render() // TODO fix
	return s
}

func (s SecretsManagerSecretResourcePolicy) GetName() string {
	return "Secrets Manager " + s.secretId + " Resource Policy"
}

func (s SecretsManagerSecretResourcePolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SecretsManagerSecretResourcePolicy) Render() {
	policy, err := s.smClient.GetResourcePolicy(
		context.TODO(),
		&sm.GetResourcePolicyInput{
			SecretId: aws.String(s.secretId),
		},
	)
	if err != nil {
		panic(err)
	}
	s.SetText(*policy.ResourcePolicy)
}
