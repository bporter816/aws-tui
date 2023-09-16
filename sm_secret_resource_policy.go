package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type SMSecretResourcePolicy struct {
	*ui.Text
	repo       *repo.SecretsManager
	secretName string
	app        *Application
}

func NewSMSecretResourcePolicy(repo *repo.SecretsManager, secretName string, app *Application) *SMSecretResourcePolicy {
	s := &SMSecretResourcePolicy{
		Text:       ui.NewText(true, "json"),
		repo:       repo,
		secretName: secretName,
		app:        app,
	}
	return s
}

func (s SMSecretResourcePolicy) GetService() string {
	return "Secrets Manager"
}

func (s SMSecretResourcePolicy) GetLabels() []string {
	return []string{s.secretName, "Resource Policy"}
}

func (s SMSecretResourcePolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretResourcePolicy) Render() {
	policy, err := s.repo.GetResourcePolicy(s.secretName)
	if err != nil {
		panic(err)
	}
	s.SetText(policy)
}
