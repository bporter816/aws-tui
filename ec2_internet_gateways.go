package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type EC2InternetGateways struct {
	*ui.Table
	view.EC2
	repo *repo.EC2
	app  *Application
}

func NewEC2InternetGateways(repo *repo.EC2, app *Application) *EC2InternetGateways {
	e := &EC2InternetGateways{
		Table: ui.NewTable([]string{
			"NAME",
			"ID",
			"OWNER",
			"ATTACHMENTS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e EC2InternetGateways) GetLabels() []string {
	return []string{"Internet Gateways"}
}

func (e EC2InternetGateways) tagsHandler() {
	internetGatewayId, err := e.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewEC2InternetGatewayTags(e.repo, internetGatewayId, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e EC2InternetGateways) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e EC2InternetGateways) Render() {
	model, err := e.repo.ListInternetGateways()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, id, ownerId, attachments string
		if n, ok := lookupTag(v.Tags, "Name"); ok {
			name = n
		} else {
			name = "-"
		}
		if v.InternetGatewayId != nil {
			id = *v.InternetGatewayId
		}
		if v.OwnerId != nil {
			ownerId = *v.OwnerId
		}
		attachments = strconv.Itoa(len(v.Attachments))
		data = append(data, []string{
			name,
			id,
			ownerId,
			attachments,
		})
	}
	e.SetData(data)
}
