package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type SMSecretResourcePolicy struct {
	*ui.Text
	view.SecretsManager
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
