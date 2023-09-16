package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type SMSecretTags struct {
	*ui.Table
	repo       *repo.SecretsManager
	secretName string
	app        *Application
}

func NewSMSecretTags(repo *repo.SecretsManager, secretName string, app *Application) *SMSecretTags {
	s := &SMSecretTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:       repo,
		secretName: secretName,
		app:        app,
	}
	return s
}

func (s SMSecretTags) GetService() string {
	return "Secrets Manager"
}

func (s SMSecretTags) GetLabels() []string {
	return []string{s.secretName, "Tags"}
}

func (s SMSecretTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretTags) Render() {
	model, err := s.repo.ListTags(s.secretName)
	if err != nil {
		panic(err)
	}
	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	s.SetData(data)
}
