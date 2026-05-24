package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEC2Instances(repo *repo.EC2, app *Application) *TableView[model.EC2Instance] {
	return NewTableView(
		"EC2",
		[]string{"Instances"},
		[]string{"NAME", "INSTANCE ID", "STATE", "PUBLIC IP", "INSTANCE TYPE", "SUBNET ID", "KEY NAME"},
		repo.ListInstances,
		func(v model.EC2Instance) []string {
			var name, state string
			if n, ok := utils.LookupEC2Tag(v.Tags, "Name"); ok {
				name = n
			}
			if v.State != nil {
				state = utils.TitleCase(string(v.State.Name))
			}
			return []string{
				name,
				utils.DerefString(v.InstanceId, ""),
				state,
				utils.DerefString(v.PublicIpAddress, ""),
				string(v.InstanceType),
				utils.DerefString(v.SubnetId, ""),
				utils.DerefString(v.KeyName, ""),
			}
		},
		func(v *TableView[model.EC2Instance]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						instanceId, err := v.GetColSelection("INSTANCE ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), instanceId, app))
					},
				},
			}
		},
	)
}
