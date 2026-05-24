package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEC2ReservedInstances(repo *repo.EC2, app *Application) *TableView[model.EC2ReservedInstance] {
	return NewTableView(
		"EC2",
		[]string{"Reserved Instances"},
		[]string{"ID", "INSTANCE TYPE", "SCOPE", "STATE"},
		func() ([]model.EC2ReservedInstance, error) {
			return repo.ListReservedInstances(nil)
		},
		func(v model.EC2ReservedInstance) []string {
			return []string{
				utils.DerefString(v.ReservedInstancesId, ""),
				string(v.InstanceType),
				string(v.Scope),
				utils.AutoCase(string(v.State)),
			}
		},
		func(v *TableView[model.EC2ReservedInstance]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						id, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), id, app))
					},
				},
			}
		},
	)
}
