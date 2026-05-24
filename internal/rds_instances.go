package internal

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	rdsTypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRDSInstances(repo *repo.RDS, app *Application, dbClusterId string) *TableView[model.RDSInstance] {
	return NewTableView(
		"RDS",
		[]string{dbClusterId, "Instances"},
		[]string{"NAME", "STATUS", "CLASS", "ENDPOINT"},
		func() ([]model.RDSInstance, error) {
			return repo.ListInstances([]rdsTypes.Filter{
				{
					Name:   aws.String("db-cluster-id"),
					Values: []string{dbClusterId},
				},
			})
		},
		func(v model.RDSInstance) []string {
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
			return []string{
				utils.DerefString(v.DBInstanceIdentifier, ""),
				status,
				utils.DerefString(v.DBInstanceClass, ""),
				endpoint,
			}
		},
		func(v *TableView[model.RDSInstance]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.DBInstanceArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.DBInstanceArn, app))
					},
				},
			}
		},
	)
}
