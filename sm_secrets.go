package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type SMSecrets struct {
	*Table
	smClient *sm.Client
	app      *Application
}

func NewSMSecrets(smClient *sm.Client, app *Application) *SMSecrets {
	s := &SMSecrets{
		Table: NewTable([]string{
			"NAME",
			"PRIMARY REGION",
			"ROTATION",
			"DESCRIPTION",
		}, 1, 0),
		smClient: smClient,
		app:      app,
	}
	return s
}

func (s SMSecrets) GetService() string {
	return "Secrets Manager"
}

func (s SMSecrets) GetLabels() []string {
	return []string{"Secrets"}
}

func (s SMSecrets) resourcePolicyHandler() {
	secretId, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	resourcePolicyView := NewSMSecretResourcePolicy(s.smClient, secretId)
	s.app.AddAndSwitch(resourcePolicyView)
}

func (s SMSecrets) tagsHandler() {
	secretId, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewSMSecretTags(s.smClient, secretId, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s SMSecrets) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Resource Policy",
			Action:      s.resourcePolicyHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s SMSecrets) Render() {
	pg := sm.NewListSecretsPaginator(
		s.smClient,
		&sm.ListSecretsInput{
			IncludePlannedDeletion: aws.Bool(true),
		},
	)
	var secrets []smTypes.SecretListEntry
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		secrets = append(secrets, out.SecretList...)
	}

	var data [][]string
	for _, v := range secrets {
		var name, primaryRegion, desc string
		var rotationEnabled bool
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
