package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type ACMPCACertificateAuthorities struct {
	*ui.Table
	view.ACMPCA
	repo  *repo.ACMPCA
	app   *Application
	model []model.ACMPCACertificateAuthority
}

func NewACMPCACertificateAuthorities(repo *repo.ACMPCA, app *Application) *ACMPCACertificateAuthorities {
	a := &ACMPCACertificateAuthorities{
		Table: ui.NewTable([]string{
			"COMMON NAME",
			"ID",
			"STATUS",
			"TYPE",
			"USAGE MODE",
			"KEY ALGO",
			"SIGNING ALGO",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return a
}

func (a ACMPCACertificateAuthorities) GetLabels() []string {
	return []string{"Certificate Authorities"}
}

func (a ACMPCACertificateAuthorities) tagsHandler() {
	row, err := a.GetRowSelection()
	if err != nil {
		return
	}
	if aa := a.model[row-1].Arn; aa != nil {
		tagsView := NewTags(a.repo, a.GetService(), *aa, a.app)
		a.app.AddAndSwitch(tagsView)
	}
}

func (a ACMPCACertificateAuthorities) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      a.tagsHandler,
		},
	}
}

func (a *ACMPCACertificateAuthorities) Render() {
	model, err := a.repo.ListCertificateAuthorities()
	if err != nil {
		panic(err)
	}
	a.model = model

	var data [][]string
	for _, v := range model {
		var commonName, id, status, caType, usageMode, keyAlgo, signingAlgo string
		if c := v.CertificateAuthorityConfiguration; c != nil {
			keyAlgo = utils.AutoCase(string(c.KeyAlgorithm))
			signingAlgo = utils.AutoCase(string(c.SigningAlgorithm))
			if s := c.Subject; s != nil {
				if s.CommonName != nil {
					commonName = *s.CommonName
				}
			}
		}
		if v.Arn != nil {
			if aa, err := arn.Parse(*v.Arn); err == nil {
				id = utils.GetResourceNameFromArn(aa)
			}
		}
		status = utils.AutoCase(string(v.Status))
		caType = utils.AutoCase(string(v.Type))
		usageMode = utils.AutoCase(string(v.UsageMode))
		data = append(data, []string{
			commonName,
			id,
			status,
			caType,
			usageMode,
			keyAlgo,
			signingAlgo,
		})
	}
	a.SetData(data)
}
