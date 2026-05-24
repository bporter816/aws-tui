package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewSMSecrets(repo *repo.SecretsManager, app *Application) *TableView[model.SecretsManagerSecret] {
	return NewTableView(
		"Secrets Manager",
		[]string{"Secrets"},
		[]string{"NAME", "PRIMARY REGION", "ROTATION", "DESCRIPTION"},
		repo.ListSecrets,
		func(v model.SecretsManagerSecret) []string {
			var rotationEnabled string
			if v.RotationEnabled != nil {
				rotationEnabled = utils.BoolToString(*v.RotationEnabled, "Yes", "No")
			}
			return []string{
				utils.DerefString(v.Name, ""),
				utils.DerefString(v.PrimaryRegion, ""),
				rotationEnabled,
				utils.DerefString(v.Description, ""),
			}
		},
		func(v *TableView[model.SecretsManagerSecret]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Resource Policy",
					Action: func() {
						secretName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSMSecretResourcePolicy(repo, secretName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'v', tcell.ModNone),
					Description: "Get Secret Value",
					Action: func() {
						secretName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSMSecretValue(repo, secretName, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						secretName, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), secretName, app))
					},
				},
			}
		},
	)
}

func init() {
	register("Secrets Manager.Secrets", func(r Repos, app *Application) Component {
		return NewSMSecrets(r.SecretsManager, app)
	})
}
