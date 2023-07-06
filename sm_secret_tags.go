package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SMSecretTags struct {
	*Table
	smClient *sm.Client
	secretId string
	app      *Application
}

func NewSMSecretTags(smClient *sm.Client, secretId string, app *Application) *SMSecretTags {
	s := &SMSecretTags{
		Table: NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		smClient: smClient,
		secretId: secretId,
		app:      app,
	}
	return s
}

func (s SMSecretTags) GetName() string {
	return fmt.Sprintf("Secrets Manager | %v | Tags", s.secretId)
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
