package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEC2KeyPairs(repo *repo.EC2, app *Application) *TableView[model.EC2KeyPair] {
	return NewTableView(
		"EC2",
		[]string{"Key Pairs"},
		[]string{"NAME", "TYPE", "FINGERPRINT", "CREATED", "ID"},
		repo.ListKeyPairs,
		func(v model.EC2KeyPair) []string {
			var created string
			if v.CreateTime != nil {
				created = v.CreateTime.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.KeyName, ""),
				string(v.KeyType),
				utils.DerefString(v.KeyFingerprint, ""),
				created,
				utils.DerefString(v.KeyPairId, ""),
			}
		},
		func(v *TableView[model.EC2KeyPair]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
					Description: "Public Key",
					Action: func() {
						keyId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewEC2KeyPairPubKey(repo, keyId, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						keyId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), keyId, app))
					},
				},
			}
		},
	)
}
