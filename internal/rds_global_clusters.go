package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
)

func NewRDSGlobalClusters(repo *repo.RDS, app *Application) *TableView[model.RDSGlobalCluster] {
	return NewTableView(
		"RDS",
		[]string{"Global Clusters"},
		[]string{"NAME", "ENGINE", "STATUS"},
		repo.ListGlobalClusters,
		func(v model.RDSGlobalCluster) []string {
			var engine, status string
			if v.Engine != nil {
				engine = *v.Engine
				if v.EngineVersion != nil {
					engine += " " + *v.EngineVersion
				}
			}
			if v.Status != nil {
				status = utils.AutoCase(*v.Status)
			}
			return []string{
				utils.DerefString(v.GlobalClusterIdentifier, ""),
				engine,
				status,
			}
		},
		nil,
	)
}
