package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type SMSecrets struct {
	*ui.Table
	view.SecretsManager
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

func (s SMSecrets) valueHandler() {
	secretName, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	valueView := NewSMSecretValue(s.repo, secretName, s.app)
	s.app.AddAndSwitch(valueView)
}

func (s SMSecrets) tagsHandler() {
	secretName, err := s.GetColSelection("NAME")
	if err != nil {
		return
	}
	tagsView := NewTags(s.repo, s.GetService(), secretName, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s SMSecrets) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Resource Policy",
			Action:      s.resourcePolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'v', tcell.ModNone),
			Description: "Get Secret Value",
			Action:      s.valueHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var rotationEnabled string
		if v.RotationEnabled != nil {
			rotationEnabled = utils.BoolToString(*v.RotationEnabled, "Yes", "No")
		}
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			utils.DerefString(v.PrimaryRegion, ""),
			rotationEnabled,
			utils.DerefString(v.Description, ""),
		})
	}
	s.SetData(data)
}
