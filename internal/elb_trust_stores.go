package internal

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewELBTrustStores(repo *repo.ELB, app *Application) *TableView[model.ELBTrustStore] {
	return NewTableView(
		"ELB",
		[]string{"Trust Stores"},
		[]string{"NAME", "STATUS", "CA CERTS", "REVOKED CERTS"},
		repo.ListTrustStores,
		func(v model.ELBTrustStore) []string {
			var caCerts, revokedCerts string
			if v.NumberOfCaCertificates != nil {
				caCerts = strconv.Itoa(int(*v.NumberOfCaCertificates))
			}
			if v.TotalRevokedEntries != nil {
				revokedCerts = strconv.FormatInt(*v.TotalRevokedEntries, 10)
			}
			return []string{
				utils.DerefString(v.Name, ""),
				utils.AutoCase(string(v.Status)),
				caCerts,
				revokedCerts,
			}
		},
		func(v *TableView[model.ELBTrustStore]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Associations",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TrustStoreArn == nil {
							return
						}
						a, err := arn.Parse(*item.TrustStoreArn)
						if err != nil {
							return
						}
						app.AddAndSwitch(NewELBTrustStoreAssociations(repo, a, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'c', tcell.ModNone),
					Description: "Certificate Bundle",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TrustStoreArn == nil {
							return
						}
						app.AddAndSwitch(NewELBTrustStoreBundle(repo, *item.TrustStoreArn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TrustStoreArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.TrustStoreArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ELB.Trust Stores", func(r Repos, app *Application) Component {
		return NewELBTrustStores(r.ELB, app)
	})
}
