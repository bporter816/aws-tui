package main

import (
	"fmt"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type KmsKeys struct {
	*ui.Table
	repo *repo.KMS
	app  *Application
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

func (k KmsKeys) GetService() string {
	return "KMS"
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
	keyId, err := k.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewKmsKeyTags(k.repo, keyId, k.app)
	k.app.AddAndSwitch(tagsView)
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
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      k.tagsHandler,
		},
	}
}

func (k KmsKeys) Render() {
	model, err := k.repo.ListKeys()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var keyId, regionality, description string
		if v.KeyId != nil {
			keyId = *v.KeyId
		}
		if v.MultiRegion != nil && *v.MultiRegion && v.MultiRegionConfiguration != nil {
			regionality = string(v.MultiRegionConfiguration.MultiRegionKeyType)
		}
		if v.Description != nil {
			description = *v.Description
		}
		data = append(data, []string{
			keyId,
			renderAliases(v.Aliases),
			utils.BoolToString(v.Enabled, "Yes", "No"),
			string(v.KeyState),
			string(v.KeySpec),
			string(v.KeyUsage),
			regionality,
			description,
		})
	}
	k.SetData(data)
}

func renderAliases(a []string) string {
	if len(a) == 0 {
		return "-"
	}
	if len(a) == 1 {
		return a[0]
	}
	return fmt.Sprintf("%v + %v more", a[0], len(a)-1)
}
