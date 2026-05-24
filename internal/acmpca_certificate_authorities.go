package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewACMPCACertificateAuthorities(repo *repo.ACMPCA, app *Application) *TableView[model.ACMPCACertificateAuthority] {
	return NewTableView(
		"ACM PCA",
		[]string{"Certificate Authorities"},
		[]string{"COMMON NAME", "ID", "STATUS", "TYPE", "USAGE MODE", "KEY ALGO", "SIGNING ALGO"},
		repo.ListCertificateAuthorities,
		func(v model.ACMPCACertificateAuthority) []string {
			var commonName, id, keyAlgo, signingAlgo string
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
			return []string{
				commonName,
				id,
				utils.AutoCase(string(v.Status)),
				utils.AutoCase(string(v.Type)),
				utils.AutoCase(string(v.UsageMode)),
				keyAlgo,
				signingAlgo,
			}
		},
		func(v *TableView[model.ACMPCACertificateAuthority]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.Arn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.Arn, app))
					},
				},
			}
		},
	)
}
