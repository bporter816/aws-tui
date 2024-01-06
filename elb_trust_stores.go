package main

import (
	"strconv"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type ELBTrustStores struct {
	*ui.Table
	repo *repo.ELB
	app  *Application
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

func (e ELBTrustStores) GetService() string {
	return "ELB"
}

func (e ELBTrustStores) GetLabels() []string {
	return []string{"Trust Stores"}
}

func (e ELBTrustStores) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (e ELBTrustStores) Render() {
	model, err := e.repo.ListTrustStores()
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
