package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type EC2KeyPairs struct {
	*ui.Table
	repo      *repo.EC2
	ec2Client *ec2.Client
	app       *Application
}

func NewEC2KeyPairs(repo *repo.EC2, ec2Client *ec2.Client, app *Application) *EC2KeyPairs {
	e := &EC2KeyPairs{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"FINGERPRINT",
			"CREATED",
			"ID",
		}, 1, 0),
		repo:      repo,
		ec2Client: ec2Client,
		app:       app,
	}
	return e
}

func (e EC2KeyPairs) GetService() string {
	return "EC2"
}

func (e EC2KeyPairs) GetLabels() []string {
	return []string{"Key Pairs"}
}

func (e EC2KeyPairs) pubKeyHandler() {
	keyId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	pubKeyView := NewEC2KeyPairPubKey(e.ec2Client, keyId, e.app)
	e.app.AddAndSwitch(pubKeyView)
}

func (e EC2KeyPairs) tagsHandler() {
	keyId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2KeyPairTags(e.ec2Client, keyId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2KeyPairs) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone),
			Description: "Public Key",
			Action:      e.pubKeyHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		var name, keyType, fingerprint, created, id string
		if v.KeyName != nil {
			name = *v.KeyName
		}
		keyType = string(v.KeyType)
		if v.KeyFingerprint != nil {
			fingerprint = *v.KeyFingerprint
		}
		if v.CreateTime != nil {
			created = v.CreateTime.Format(utils.DefaultTimeFormat)
		}
		if v.KeyPairId != nil {
			id = *v.KeyPairId
		}
		data = append(data, []string{
			name,
			keyType,
			fingerprint,
			created,
			id,
		})
	}
	e.SetData(data)
}
