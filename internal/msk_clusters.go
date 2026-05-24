package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewMSKClusters(repo *repo.MSK, app *Application) *TableView[model.MSKCluster] {
	return NewTableView(
		"MSK",
		[]string{"Clusters"},
		[]string{"NAME", "TYPE", "STATE"},
		repo.ListClusters,
		func(v model.MSKCluster) []string {
			return []string{
				utils.DerefString(v.ClusterName, ""),
				utils.AutoCase(string(v.ClusterType)),
				utils.AutoCase(string(v.State)),
				// TODO add provisioned and serverless details
			}
		},
		func(v *TableView[model.MSKCluster]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ClusterArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ClusterArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("MSK.Clusters", func(r Repos, app *Application) Component {
		return NewMSKClusters(r.MSK, app)
	})
}
