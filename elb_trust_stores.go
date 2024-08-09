package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type ELBTrustStores struct {
	*ui.Table
	view.ELB
	repo  *repo.ELB
	app   *Application
	model []model.ELBTrustStore
}

func NewELBTrustStores(repo *repo.ELB, app *Application) *ELBTrustStores {
	e := &ELBTrustStores{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"CA CERTS",
			"REVOKED CERTS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ELBTrustStores) GetLabels() []string {
	return []string{"Trust Stores"}
}

func (e ELBTrustStores) associationsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	a, err := arn.Parse(*e.model[row-1].TrustStoreArn)
	if err != nil {
		return
	}
	associationsView := NewELBTrustStoreAssociations(e.repo, a, e.app)
	e.app.AddAndSwitch(associationsView)
}

func (e ELBTrustStores) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}

	tagsView := NewTags(e.repo, e.GetService(), *e.model[row-1].TrustStoreArn, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ELBTrustStores) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Associations",
			Action:      e.associationsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBTrustStores) Render() {
	model, err := e.repo.ListTrustStores()
	e.model = model
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, caCerts, revokedCerts string
		if v.Name != nil {
			name = *v.Name
		}
		if v.NumberOfCaCertificates != nil {
			caCerts = strconv.Itoa(int(*v.NumberOfCaCertificates))
		}
		if v.TotalRevokedEntries != nil {
			revokedCerts = strconv.FormatInt(*v.TotalRevokedEntries, 10)
		}
		data = append(data, []string{
			name,
			utils.AutoCase(string(v.Status)),
			caCerts,
			revokedCerts,
		})
	}
	e.SetData(data)
}
