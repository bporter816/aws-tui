package internal

import (
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/view"
)

type EC2KeyPairPubKey struct {
	*ui.Text
	view.EC2
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
