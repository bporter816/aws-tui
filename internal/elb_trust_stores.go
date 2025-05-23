package internal

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
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

func (e ELBTrustStores) bundleHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if a := e.model[row-1].TrustStoreArn; a != nil {
		certView := NewELBTrustStoreBundle(e.repo, *a, e.app)
		e.app.AddAndSwitch(certView)
	}
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
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Certificate Bundle",
			Action:      e.bundleHandler,
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
		var caCerts, revokedCerts string
		if v.NumberOfCaCertificates != nil {
			caCerts = strconv.Itoa(int(*v.NumberOfCaCertificates))
		}
		if v.TotalRevokedEntries != nil {
			revokedCerts = strconv.FormatInt(*v.TotalRevokedEntries, 10)
		}
		data = append(data, []string{
			utils.DerefString(v.Name, ""),
			utils.AutoCase(string(v.Status)),
			caCerts,
			revokedCerts,
		})
	}
	e.SetData(data)
}
