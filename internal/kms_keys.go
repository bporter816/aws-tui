package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewKmsKeys(repo *repo.KMS, app *Application) *TableView[model.KMSKey] {
	return NewTableView(
		"KMS",
		[]string{"Keys"},
		[]string{"ID", "ALIASES", "ENABLED", "STATE", "SPEC", "USAGE", "REGIONALITY", "DESCRIPTION"},
		repo.ListKeys,
		func(v model.KMSKey) []string {
			var regionality string
			if v.MultiRegion != nil && *v.MultiRegion && v.MultiRegionConfiguration != nil {
				regionality = string(v.MultiRegionConfiguration.MultiRegionKeyType)
			}
			return []string{
				utils.DerefString(v.KeyId, ""),
				utils.FormatKMSAliases(v.Aliases),
				utils.BoolToString(v.Enabled, "Yes", "No"),
				string(v.KeyState),
				string(v.KeySpec),
				string(v.KeyUsage),
				regionality,
				utils.DerefString(v.Description, ""),
			}
		},
		func(v *TableView[model.KMSKey]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Key Policy",
					Action: func() {
						keyId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewKmsKeyPolicy(repo, keyId, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
					Description: "Grants",
					Action: func() {
						keyId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewKmsKeyGrants(repo, keyId, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.Arn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.Arn, app))
					},
				},
			}
		},
	)
}
