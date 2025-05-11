package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ACMCertificates struct {
	*ui.Table
	view.ACM
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
		tagsView := NewTags(a.repo, a.GetService(), *arn, a.app)
		a.app.AddAndSwitch(tagsView)
	}
}

func (a ACMCertificates) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
			Description: "Certificate",
			Action:      a.certificateHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var id, inUse string
		if v.CertificateArn != nil {
			arn, err := arn.Parse(*v.CertificateArn)
			if err == nil {
				id = utils.GetResourceNameFromArn(arn)
			}
		}
		if v.InUse != nil {
			inUse = utils.BoolToString(*v.InUse, "Yes", "No")
		}
		data = append(data, []string{
			id,
			utils.DerefString(v.DomainName, ""),
			utils.AutoCase(string(v.Type)),
			string(v.KeyAlgorithm),
			utils.AutoCase(string(v.Status)),
			inUse,
			string(v.RenewalEligibility),
		})
	}
	a.SetData(data)
}
