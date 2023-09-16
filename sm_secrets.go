package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type SMSecrets struct {
	*ui.Table
	repo *repo.SecretsManager
	app  *Application
}

func NewSMSecrets(repo *repo.SecretsManager, app *Application) *SMSecrets {
	s := &SMSecrets{
		Table: ui.NewTable([]string{
			"NAME",
			"PRIMARY REGION",
			"ROTATION",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
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
	secretName, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	resourcePolicyView := NewSMSecretResourcePolicy(s.repo, secretName, s.app)
	s.app.AddAndSwitch(resourcePolicyView)
}

func (s SMSecrets) tagsHandler() {
	secretName, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewSMSecretTags(s.repo, secretName, s.app)
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
	model, err := s.repo.ListSecrets()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
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
