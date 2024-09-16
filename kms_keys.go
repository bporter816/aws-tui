package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type KmsKeys struct {
	*ui.Table
	view.KMS
	repo  *repo.KMS
	app   *Application
	model []model.KMSKey
}

func NewKmsKeys(repo *repo.KMS, app *Application) *KmsKeys {
	k := &KmsKeys{
		Table: ui.NewTable([]string{
			"ID",
			"ALIASES",
			"ENABLED",
			"STATE",
			"SPEC",
			"USAGE",
			"REGIONALITY",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return k
}

func (k KmsKeys) GetLabels() []string {
	return []string{"Keys"}
}

func (k KmsKeys) keyPolicyHandler() {
	keyId, err := k.GetColSelection("ID")
	if err != nil {
		return
	}
	policyView := NewKmsKeyPolicy(k.repo, keyId, k.app)
	k.app.AddAndSwitch(policyView)
}

func (k KmsKeys) grantsHandler() {
	keyId, err := k.GetColSelection("ID")
	if err != nil {
		return
	}
	grantsView := NewKmsKeyGrants(k.repo, keyId, k.app)
	k.app.AddAndSwitch(grantsView)
}

func (k KmsKeys) tagsHandler() {
	row, err := k.GetRowSelection()
	if err != nil {
		return
	}
	if arn := k.model[row-1].Arn; arn != nil {
		tagsView := NewTags(k.repo, k.GetService(), *arn, k.app)
		k.app.AddAndSwitch(tagsView)
	}
}

func (k KmsKeys) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Key Policy",
			Action:      k.keyPolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone),
			Description: "Grants",
			Action:      k.grantsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      k.tagsHandler,
		},
	}
}

func (k *KmsKeys) Render() {
	model, err := k.repo.ListKeys()
	k.model = model
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var regionality string
		if v.MultiRegion != nil && *v.MultiRegion && v.MultiRegionConfiguration != nil {
			regionality = string(v.MultiRegionConfiguration.MultiRegionKeyType)
		}
		data = append(data, []string{
			utils.DerefString(v.KeyId, ""),
			utils.FormatKMSAliases(v.Aliases),
			utils.BoolToString(v.Enabled, "Yes", "No"),
			string(v.KeyState),
			string(v.KeySpec),
			string(v.KeyUsage),
			regionality,
			utils.DerefString(v.Description, ""),
		})
	}
	k.SetData(data)
}
