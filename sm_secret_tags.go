package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bporter816/aws-tui/ui"
)

type SMSecretTags struct {
	*ui.Table
	smClient *sm.Client
	secretId string
	app      *Application
}

func NewSMSecretTags(smClient *sm.Client, secretId string, app *Application) *SMSecretTags {
	s := &SMSecretTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		smClient: smClient,
		secretId: secretId,
		app:      app,
	}
	return s
}

func (s SMSecretTags) GetService() string {
	return "Secrets Manager"
}

func (s SMSecretTags) GetLabels() []string {
	return []string{s.secretId, "Tags"}
}

func (s SMSecretTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretTags) Render() {
	out, err := s.smClient.DescribeSecret(
		context.TODO(),
		&sm.DescribeSecretInput{
			SecretId: aws.String(s.secretId),
		},
	)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range out.Tags {
		data = append(data, []string{
			*v.Key,
			*v.Value,
		})
	}
	s.SetData(data)
}
