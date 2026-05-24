package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewACMCertificates(repo *repo.ACM, app *Application) *TableView[model.ACMCertificate] {
	return NewTableView(
		"ACM",
		[]string{"Certificates"},
		[]string{"ID", "DOMAIN", "TYPE", "ALGORITHM", "STATUS", "IN USE", "RENEWAL ELIGIBILITY"},
		repo.ListCertificates,
		func(v model.ACMCertificate) []string {
			var id, inUse string
			if v.CertificateArn != nil {
				a, err := arn.Parse(*v.CertificateArn)
				if err == nil {
					id = utils.GetResourceNameFromArn(a)
				}
			}
			if v.InUse != nil {
				inUse = utils.BoolToString(*v.InUse, "Yes", "No")
			}
			return []string{
				id,
				utils.DerefString(v.DomainName, ""),
				utils.AutoCase(string(v.Type)),
				string(v.KeyAlgorithm),
				utils.AutoCase(string(v.Status)),
				inUse,
				string(v.RenewalEligibility),
			}
		},
		func(v *TableView[model.ACMCertificate]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
					Description: "Certificate",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.CertificateArn == nil {
							return
						}
						app.AddAndSwitch(NewACMCertificateDetails(repo, "Certificate", *item.CertificateArn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.CertificateArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.CertificateArn, app))
					},
				},
			}
		},
	)
}
