// TODO add remaining days based on start time and duration

package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRDSReservedInstances(repo *repo.RDS, app *Application) *TableView[model.RDSReservedInstance] {
	return NewTableView(
		"RDS",
		[]string{"Reserved Instances"},
		[]string{"RESERVATION ID", "LEASE ID", "PRODUCT", "OFFERING TYPE", "CLASS", "COUNT", "STATUS", "MULTI AZ", "START TIME"},
		repo.ListReservedInstances,
		func(v model.RDSReservedInstance) []string {
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
			return []string{
				utils.DerefString(v.ReservedDBInstanceId, ""),
				utils.DerefString(v.LeaseId, ""),
				utils.DerefString(v.ProductDescription, ""),
				utils.DerefString(v.OfferingType, ""),
				utils.DerefString(v.DBInstanceClass, ""),
				count,
				status,
				multiAZ,
				startTime,
			}
		},
		func(v *TableView[model.RDSReservedInstance]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.ReservedDBInstanceArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.ReservedDBInstanceArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("RDS.Reserved Instances", func(r Repos, app *Application) Component {
		return NewRDSReservedInstances(r.RDS, app)
	})
}
