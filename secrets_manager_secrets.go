package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type SecretsManagerSecrets struct {
	*Table
	smClient *sm.Client
	app      *Application
}

func NewSecretsManagerSecrets(smClient *sm.Client, app *Application) *SecretsManagerSecrets {
	s := &SecretsManagerSecrets{
		Table: NewTable([]string{
			"NAME",
			"PRIMARY REGION",
			"ROTATION",
			"DESCRIPTION",
		}, 1, 0),
		smClient: smClient,
		app:      app,
	}
	s.Render() // TODO fix
	return s
}

func (s SecretsManagerSecrets) GetName() string {
	return "Secrets Manager"
}

func (s SecretsManagerSecrets) resourcePolicyHandler() {
	r, _ := s.GetSelection()
	secretId := s.GetCell(r, 0).Text
	resourcePolicyView := NewSecretsManagerSecretResourcePolicy(s.smClient, secretId)
	s.app.AddAndSwitch("sm.resourcepolicy", resourcePolicyView)
}

func (s SecretsManagerSecrets) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "",
			Action:      s.resourcePolicyHandler,
		},
	}
}

func (s SecretsManagerSecrets) Render() {
	secretsPaginator := sm.NewListSecretsPaginator(s.smClient, &sm.ListSecretsInput{IncludePlannedDeletion: aws.Bool(true)})
	var secrets []smTypes.SecretListEntry
	for secretsPaginator.HasMorePages() {
		out, err := secretsPaginator.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		secrets = append(secrets, out.SecretList...)
	}

	var data [][]string
	for _, v := range secrets {
		var name string
		var primaryRegion string
		var rotationEnabled bool
		var desc string
		if v.Name != nil {
			name = *v.Name
		}
		if v.PrimaryRegion != nil {
			primaryRegion = *v.PrimaryRegion
		}
		if v.RotationEnabled != nil {
			rotationEnabled = *v.RotationEnabled
		}
		if v.Description != nil {
			desc = *v.Description
		}
		data = append(data, []string{
			name,
			primaryRegion,
			strconv.FormatBool(rotationEnabled),
			desc,
		})
	}
	s.SetData(data)
}
