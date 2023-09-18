package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type EC2KeyPairPubKey struct {
	*ui.Text
	repo      *repo.EC2
	keyPairId string
	app       *Application
}

func NewEC2KeyPairPubKey(repo *repo.EC2, keyPairId string, app *Application) *EC2KeyPairPubKey {
	e := &EC2KeyPairPubKey{
		Text:      ui.NewText(false, ""),
		repo:      repo,
		keyPairId: keyPairId,
		app:       app,
	}
	return e
}

func (e EC2KeyPairPubKey) GetService() string {
	return "EC2"
}

func (e EC2KeyPairPubKey) GetLabels() []string {
	return []string{e.keyPairId, "Public Key"}
}

func (e EC2KeyPairPubKey) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2KeyPairPubKey) Render() {
	key, err := e.repo.GetPublicKey(e.keyPairId)
	if err != nil {
		panic(err)
	}

	e.SetText(key)
}
