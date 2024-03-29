package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type RDSParameterGroups struct {
	*ui.Table
	view.RDS
	repo *repo.RDS
	app  *Application
}

func NewRDSParameterGroups(repo *repo.RDS, app *Application) *RDSParameterGroups {
	r := &RDSParameterGroups{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"FAMILY",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r RDSParameterGroups) GetLabels() []string {
	return []string{"Parameter Groups"}
}

func (r RDSParameterGroups) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r RDSParameterGroups) Render() {
	clusterParameterGroups, err := r.repo.ListClusterParameterGroups()
	if err != nil {
		panic(err)
	}

	instanceParameterGroups, err := r.repo.ListInstanceParameterGroups()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range clusterParameterGroups {
		var name, family, desc string
		if v.DBClusterParameterGroupName != nil {
			name = *v.DBClusterParameterGroupName
		}
		if v.DBParameterGroupFamily != nil {
			family = *v.DBParameterGroupFamily
		}
		if v.Description != nil {
			desc = *v.Description
		}
		data = append(data, []string{
			name,
			"Cluster",
			family,
			desc,
		})
	}
	for _, v := range instanceParameterGroups {
		var name, family, desc string
		if v.DBParameterGroupName != nil {
			name = *v.DBParameterGroupName
		}
		if v.DBParameterGroupFamily != nil {
			family = *v.DBParameterGroupFamily
		}
		if v.Description != nil {
			desc = *v.Description
		}
		data = append(data, []string{
			name,
			"Instance",
			family,
			desc,
		})
	}
	r.SetData(data)
}
