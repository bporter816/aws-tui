// TODO add remaining days based on start time and duration

package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

type RDSReservedInstances struct {
	*ui.Table
	view.RDS
	repo  *repo.RDS
	app   *Application
	model []model.RDSReservedInstance
}

func NewRDSReservedInstances(repo *repo.RDS, app *Application) *RDSReservedInstances {
	r := &RDSReservedInstances{
		Table: ui.NewTable([]string{
			"RESERVATION ID",
			"LEASE ID",
			"PRODUCT",
			"OFFERING TYPE",
			"CLASS",
			"COUNT",
			"STATUS",
			"MULTI AZ",
			"START TIME",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r RDSReservedInstances) GetLabels() []string {
	return []string{"Reserved Instances"}
}

func (r RDSReservedInstances) tagsHandler() {
	row, err := r.GetRowSelection()
	if err != nil || r.model[row-1].ReservedDBInstanceArn == nil {
		return
	}
	tagsView := NewTags(r.repo, r.GetService(), *r.model[row-1].ReservedDBInstanceArn, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r RDSReservedInstances) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r *RDSReservedInstances) Render() {
	model, err := r.repo.ListReservedInstances()
	if err != nil {
		panic(err)
	}
	r.model = model

	var data [][]string
	for _, v := range model {
		var count, status, multiAZ, startTime string
		if v.DBInstanceCount != nil {
			count = strconv.Itoa(int(*v.DBInstanceCount))
		}
		if v.State != nil {
			status = utils.AutoCase(*v.State)
		}
		if v.MultiAZ != nil {
			multiAZ = utils.BoolToString(*v.MultiAZ, "Yes", "No")
		}
		if v.StartTime != nil {
			startTime = v.StartTime.Format(utils.DefaultTimeFormat)
		}
		data = append(data, []string{
			utils.DerefString(v.ReservedDBInstanceId, ""),
			utils.DerefString(v.LeaseId, ""),
			utils.DerefString(v.ProductDescription, ""),
			utils.DerefString(v.OfferingType, ""),
			utils.DerefString(v.DBInstanceClass, ""),
			count,
			status,
			multiAZ,
			startTime,
		})
	}
	r.SetData(data)
}
