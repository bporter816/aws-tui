package internal

import (
	"strconv"

	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRDSClusters(repo *repo.RDS, app *Application) *TableView[model.RDSCluster] {
	return NewTableView(
		"RDS",
		[]string{"Clusters"},
		[]string{"NAME", "PORT", "INSTANCES", "ENGINE", "PARAM GROUP", "CUSTOM ENDPOINTS"},
		func() ([]model.RDSCluster, error) {
			return repo.ListClusters([]rdsTypes.Filter{})
		},
		func(v model.RDSCluster) []string {
			var port, engine string
			if v.Port != nil {
				port = strconv.Itoa(int(*v.Port))
			}
			if v.Engine != nil {
				engine = *v.Engine
				if v.EngineVersion != nil {
					engine += " " + *v.EngineVersion
				}
				if v.EngineMode != nil {
					engine += " (" + *v.EngineMode + ")"
				}
			}
			return []string{
				utils.DerefString(v.DBClusterIdentifier, ""),
				port,
				strconv.Itoa(len(v.DBClusterMembers)),
				engine,
				utils.DerefString(v.DBClusterParameterGroup, ""),
				strconv.Itoa(len(v.CustomEndpoints)),
			}
		},
		func(v *TableView[model.RDSCluster]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
					Description: "Instances",
					Action: func() {
						clusterId, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewRDSInstances(repo, app, clusterId))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
					Description: "Endpoints",
					Action: func() {
						clusterId, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewRDSEndpoints(repo, app, clusterId))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.DBClusterArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.DBClusterArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("RDS.Clusters", func(r Repos, app *Application) Component {
		return NewRDSClusters(r.RDS, app)
	})
}
