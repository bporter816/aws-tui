package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ACMCertificates struct {
	*ui.Table
	repo  *repo.ACM
	app   *Application
	model []model.ACMCertificate
}

func NewACMCertificates(repo *repo.ACM, app *Application) *ACMCertificates {
	a := &ACMCertificates{
		Table: ui.NewTable([]string{
			"ID",
			"DOMAIN",
			"TYPE",
			"ALGORITHM",
			"STATUS",
			"IN USE",
			"RENEWAL ELIGIBILITY",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return a
}

func (a ACMCertificates) GetService() string {
	return "ACM"
}

func (a ACMCertificates) GetLabels() []string {
	return []string{"Certificates"}
}

func (a ACMCertificates) certificateHandler() {
	row, err := a.GetRowSelection()
	if err != nil {
		return
	}
	if arn := a.model[row-1].CertificateArn; arn != nil {
		certificateView := NewACMCertificateDetails(a.repo, "Certificate", *arn, a.app)
		a.app.AddAndSwitch(certificateView)
	}
}

func (a ACMCertificates) tagsHandler() {
	row, err := a.GetRowSelection()
	if err != nil {
		return
	}
	if arn := a.model[row-1].CertificateArn; arn != nil {
		tagsView := NewACMTags(a.repo, *arn, a.app)
		a.app.AddAndSwitch(tagsView)
	}
}

func (a ACMCertificates) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Certificate",
			Action:      a.certificateHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      a.tagsHandler,
		},
	}
}

func (a *ACMCertificates) Render() {
	model, err := a.repo.ListCertificates()
	if err != nil {
		panic(err)
	}
	a.model = model

	var data [][]string
	for _, v := range model {
		var id, domainName, certificateType, algorithm, status, inUse, renewalEligibility string
		if v.CertificateArn != nil {
			arn, err := arn.Parse(*v.CertificateArn)
			if err == nil {
				id = utils.GetResourceNameFromArn(arn)
			}
		}
		if v.DomainName != nil {
			domainName = *v.DomainName
		}
		certificateType = utils.AutoCase(string(v.Type))
		algorithm = string(v.KeyAlgorithm)
		status = utils.AutoCase(string(v.Status))
		if v.InUse != nil {
			inUse = utils.BoolToString(*v.InUse, "Yes", "No")
		}
		renewalEligibility = string(v.RenewalEligibility)
		data = append(data, []string{
			id,
			domainName,
			certificateType,
			algorithm,
			status,
			inUse,
			renewalEligibility,
		})
	}
	a.SetData(data)
}
