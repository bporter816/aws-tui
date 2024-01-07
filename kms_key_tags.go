package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type KmsKeyTags struct {
	*ui.Table
	view.KMS
	repo  *repo.KMS
	keyId string
	app   *Application
}

func NewKmsKeyTags(repo *repo.KMS, keyId string, app *Application) *KmsKeyTags {
	k := &KmsKeyTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:  repo,
		app:   app,
		keyId: keyId,
	}
	return k
}

func (k KmsKeyTags) GetLabels() []string {
	return []string{k.keyId, "Tags"}
}

func (k KmsKeyTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyTags) Render() {
	model, err := k.repo.ListTags(k.keyId)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	k.SetData(data)
}
