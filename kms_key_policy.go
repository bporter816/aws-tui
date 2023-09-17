package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type KmsKeyPolicy struct {
	*ui.Text
	repo  *repo.KMS
	keyId string
	app   *Application
}

func NewKmsKeyPolicy(repo *repo.KMS, keyId string, app *Application) *KmsKeyPolicy {
	k := &KmsKeyPolicy{
		Text:  ui.NewText(true, "json"),
		repo:  repo,
		keyId: keyId,
		app:   app,
	}
	return k
}

func (k KmsKeyPolicy) GetService() string {
	return "KMS"
}

func (k KmsKeyPolicy) GetLabels() []string {
	return []string{k.keyId, "Key Policy"}
}

func (k KmsKeyPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (k KmsKeyPolicy) Render() {
	policy, err := k.repo.GetKeyPolicy(k.keyId)
	if err != nil {
		panic(err)
	}
	k.SetText(policy)
}
