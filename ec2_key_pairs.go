package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2KeyPairs struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2KeyPairs(repo *repo.EC2, app *Application) *EC2KeyPairs {
	e := &EC2KeyPairs{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"FINGERPRINT",
			"CREATED",
			"ID",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2KeyPairs) GetLabels() []string {
	return []string{"Key Pairs"}
}

func (e EC2KeyPairs) pubKeyHandler() {
	keyId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	pubKeyView := NewEC2KeyPairPubKey(e.repo, keyId, e.app)
	e.app.AddAndSwitch(pubKeyView)
}

func (e EC2KeyPairs) tagsHandler() {
	keyId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(e.repo, e.GetService(), keyId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2KeyPairs) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Public Key",
			Action:      e.pubKeyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2KeyPairs) Render() {
	model, err := e.repo.ListKeyPairs()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var created string
		if v.CreateTime != nil {
			created = v.CreateTime.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.KeyName, ""),
			string(v.KeyType),
			utils.DerefString(v.KeyFingerprint, ""),
			created,
			utils.DerefString(v.KeyPairId, ""),
		})
	}
	e.SetData(data)
}
