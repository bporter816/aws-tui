package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEKSClusters(repo *repo.EKS, app *Application) *TableView[model.EKSCluster] {
	return NewTableView(
		"EKS",
		[]string{"Clusters"},
		[]string{"NAME", "STATUS", "K8S VERSION", "OIDC ISSUER", "CREATED"},
		repo.ListClusters,
		func(v model.EKSCluster) []string {
			var oidcIssuer, created string
			if v.Identity != nil && v.Identity.Oidc != nil && v.Identity.Oidc.Issuer != nil {
				oidcIssuer = *v.Identity.Oidc.Issuer
			}
			if v.CreatedAt != nil {
				created = v.CreatedAt.Format(utils.DefaultTimeFormat)
			}
			return []string{
				utils.DerefString(v.Name, ""),
				utils.TitleCase(string(v.Status)),
				utils.DerefString(v.Version, ""),
				oidcIssuer,
				created,
			}
		},
		func(v *TableView[model.EKSCluster]) []KeyAction {
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
