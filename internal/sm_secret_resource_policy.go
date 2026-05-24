package internal

import (
	"context"

	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
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

func (s SMSecretResourcePolicy) GetService() string { return "Secrets Manager" }

func (s SMSecretResourcePolicy) GetLabels() []string {
	return []string{s.secretName, "Resource Policy"}
}

func (s SMSecretResourcePolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretResourcePolicy) Render(ctx context.Context) error {
	policy, err := s.repo.GetResourcePolicy(ctx, s.secretName)
	if err != nil {
		return err
	}
	s.SetText(policy)
	return nil
}
