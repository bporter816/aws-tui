package main

import (
	"strconv"

	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type RDSClusters struct {
	*ui.Table
	view.RDS
	repo *repo.RDS
	app  *Application
}

func NewRDSClusters(repo *repo.RDS, app *Application) *RDSClusters {
	r := &RDSClusters{
		Table: ui.NewTable([]string{
			"NAME",
			"PORT",
			"INSTANCES",
			"ENGINE",
			"PARAM GROUP",
			"CUSTOM ENDPOINTS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r RDSClusters) GetLabels() []string {
	return []string{"Clusters"}
}

func (r RDSClusters) endpointsHandler() {
	clusterId, err := r.GetColSelection("NAME")
	if err != nil {
		return
	}
	endpointsView := NewRDSEndpoints(r.repo, r.app, clusterId)
	r.app.AddAndSwitch(endpointsView)
}

func (r RDSClusters) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'e', tcell.ModNone),
			Description: "Endpoints",
			Action:      r.endpointsHandler,
		},
	}
}

func (r RDSClusters) Render() {
	model, err := r.repo.ListClusters([]rdsTypes.Filter{})
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var name, port, instances, engine, paramGroup, customEndpoints string
		if v.DBClusterIdentifier != nil {
			name = *v.DBClusterIdentifier
		}
		if v.Port != nil {
			port = strconv.Itoa(int(*v.Port))
		}
		instances = strconv.Itoa(len(v.DBClusterMembers))
		if v.Engine != nil {
			engine = *v.Engine
			if v.EngineVersion != nil {
				engine += " " + *v.EngineVersion
			}
			if v.EngineMode != nil {
				engine += " (" + *v.EngineMode + ")"
			}
		}
		if v.DBClusterParameterGroup != nil {
			paramGroup = *v.DBClusterParameterGroup
		}
		customEndpoints = strconv.Itoa(len(v.CustomEndpoints))
		data = append(data, []string{
			name,
			port,
			instances,
			engine,
			paramGroup,
			customEndpoints,
		})
	}
	r.SetData(data)
}
