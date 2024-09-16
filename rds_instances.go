package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type RDSInstances struct {
	*ui.Table
	view.RDS
	repo        *repo.RDS
	app         *Application
	dbClusterId string
	model       []model.RDSInstance
}

func NewRDSInstances(repo *repo.RDS, app *Application, dbClusterId string) *RDSInstances {
	r := &RDSInstances{
		Table: ui.NewTable([]string{
			"NAME",
			"STATUS",
			"CLASS",
			"ENDPOINT",
		}, 1, 0),
		repo:        repo,
		app:         app,
		dbClusterId: dbClusterId,
	}
	return r
}

func (r RDSInstances) GetLabels() []string {
	return []string{r.dbClusterId, "Instances"}
}

func (r RDSInstances) tagsHandler() {
	row, err := r.GetRowSelection()
	if err != nil || r.model[row-1].DBInstanceArn == nil {
		return
	}
	tagsView := NewTags(r.repo, r.GetService(), *r.model[row-1].DBInstanceArn, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r RDSInstances) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r *RDSInstances) Render() {
	model, err := r.repo.ListInstances([]rdsTypes.Filter{
		{
			Name:   aws.String("db-cluster-id"),
			Values: []string{r.dbClusterId},
		},
	})
	if err != nil {
		panic(err)
	}
	r.model = model

	var data [][]string
	for _, v := range model {
		var status, endpoint string
		if v.DBInstanceStatus != nil {
			status = utils.AutoCase(*v.DBInstanceStatus)
		}
		if v.Endpoint != nil {
			if v.Endpoint.Address != nil {
				endpoint = *v.Endpoint.Address
				if v.Endpoint.Port != nil {
					endpoint += ":" + strconv.Itoa(int(*v.Endpoint.Port))
				}
			}
		}
		data = append(data, []string{
			utils.DerefString(v.DBInstanceIdentifier, ""),
			status,
			utils.DerefString(v.DBInstanceClass, ""),
			endpoint,
		})
	}
	r.SetData(data)
}
