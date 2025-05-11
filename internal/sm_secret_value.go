package internal

import (
	"encoding/json"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/rivo/tview"
)

type SMSecretValue struct {
	tview.Primitive
	view.SecretsManager
	secretName string
	valueKV    map[string]string
	valueText  string
	repo       *repo.SecretsManager
	app        *Application
}

func NewSMSecretValue(repo *repo.SecretsManager, secretName string, app *Application) *SMSecretValue {
	s := &SMSecretValue{
		secretName: secretName,
		repo:       repo,
		app:        app,
	}

	secretValue, err := s.repo.GetSecretValue(s.secretName)
	if err != nil {
		panic(err)
	}

	var kv map[string]string
	err = json.Unmarshal([]byte(secretValue), &kv)
	if err == nil {
		s.Primitive = ui.NewTable([]string{"KEY", "VALUE"}, 1, 0)
		s.valueKV = kv
	} else {
		s.Primitive = ui.NewText(false, "")
		s.valueText = secretValue
	}
	return s
}

func (s SMSecretValue) GetLabels() []string {
	return []string{s.secretName, "Value"}
}

func (s SMSecretValue) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SMSecretValue) Render() {
	if v, ok := s.Primitive.(*ui.Table); ok {
		var data [][]string
		for k, v := range s.valueKV {
			data = append(data, []string{k, v})
		}
		v.SetData(data)
	} else if v, ok := s.Primitive.(*ui.Text); ok {
		v.SetText(s.valueText)
	} else {
		panic("invalid secret value type")
	}
}
