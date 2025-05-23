package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
)

type RDSGlobalClusters struct {
	*ui.Table
	view.RDS
	repo  *repo.RDS
	app   *Application
	model []model.RDSGlobalCluster
}

func NewRDSGlobalClusters(repo *repo.RDS, app *Application) *RDSGlobalClusters {
	r := &RDSGlobalClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"ENGINE",
			"STATUS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r RDSGlobalClusters) GetLabels() []string {
	return []string{"Global Clusters"}
}

func (r RDSGlobalClusters) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r *RDSGlobalClusters) Render() {
	model, err := r.repo.ListGlobalClusters()
	if err != nil {
		panic(err)
	}
	r.model = model

	var data [][]string
	for _, v := range model {
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
		data = append(data, []string{
			utils.DerefString(v.GlobalClusterIdentifier, ""),
			engine,
			status,
		})
	}
	r.SetData(data)
}
