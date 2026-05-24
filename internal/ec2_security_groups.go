package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewEC2SecurityGroups(repo *repo.EC2, app *Application) *TableView[model.EC2SecurityGroup] {
	return NewTableView(
		"EC2",
		[]string{"Security Groups"},
		[]string{"NAME", "ID", "VPC ID", "INGRESS RULES", "EGRESS RULES", "DESCRIPTION"},
		repo.ListSecurityGroups,
		func(v model.EC2SecurityGroup) []string {
			return []string{
				utils.DerefString(v.GroupName, ""),
				utils.DerefString(v.GroupId, ""),
				utils.DerefString(v.VpcId, ""),
				strconv.Itoa(len(v.IpPermissions)),
				strconv.Itoa(len(v.IpPermissionsEgress)),
				utils.DerefString(v.Description, ""),
			}
		},
		func(v *TableView[model.EC2SecurityGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'u', tcell.ModNone),
					Description: "Rules",
					Action: func() {
						sgId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewEC2SecurityGroupRules(repo, sgId, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						sgId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), sgId, app))
					},
				},
			}
		},
	)
}
