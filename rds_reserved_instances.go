// TODO add remaining days based on start time and duration

package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
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

func (r RDSReservedInstances) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r *RDSReservedInstances) Render() {
	model, err := r.repo.ListReservedInstances()
	if err != nil {
		panic(err)
	}
	r.model = model

	var data [][]string
	for _, v := range model {
		var reservationId, leaseId, product, offeringType, class, count, status, multiAZ, startTime string
		if v.ReservedDBInstanceId != nil {
			reservationId = *v.ReservedDBInstanceId
		}
		if v.LeaseId != nil {
			leaseId = *v.LeaseId
		}
		if v.ProductDescription != nil {
			product = *v.ProductDescription
		}
		if v.OfferingType != nil {
			offeringType = *v.OfferingType
		}
		if v.DBInstanceClass != nil {
			class = *v.DBInstanceClass
		}
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
			reservationId,
			leaseId,
			product,
			offeringType,
			class,
			count,
			status,
			multiAZ,
			startTime,
		})
	}
	r.SetData(data)
}
