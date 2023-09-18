package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type EC2KeyPairTags struct {
	*ui.Table
	repo      *repo.EC2
	keyPairId string
	app       *Application
}

func NewEC2KeyPairTags(repo *repo.EC2, keyPairId string, app *Application) *EC2KeyPairTags {
	e := &EC2KeyPairTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:      repo,
		keyPairId: keyPairId,
		app:       app,
	}
	return e
}

func (e EC2KeyPairTags) GetService() string {
	return "EC2"
}

func (e EC2KeyPairTags) GetLabels() []string {
	return []string{e.keyPairId, "Tags"}
}

func (e EC2KeyPairTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e EC2KeyPairTags) Render() {
	model, err := e.repo.ListKeyPairTags(e.keyPairId)
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
	e.SetData(data)
}
