package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bporter816/aws-tui/ui"
)

type SMSecretResourcePolicy struct {
	*ui.Text
	smClient *sm.Client
	secretId string
	app      *Application
}

func NewSMSecretResourcePolicy(smClient *sm.Client, secretId string, app *Application) *SMSecretResourcePolicy {
	s := &SMSecretResourcePolicy{
		Text:     ui.NewText(true, "json"),
		smClient: smClient,
		secretId: secretId,
		app:      app,
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
