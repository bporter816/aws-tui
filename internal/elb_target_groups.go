package internal

import (
	"strconv"

	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewELBTargetGroups(repo *repo.ELB, app *Application) *TableView[model.ELBTargetGroup] {
	return NewTableView(
		"ELB",
		[]string{"Target Groups"},
		[]string{"NAME", "PORT", "PROTOCOL", "TARGET TYPE", "VPC"},
		repo.ListTargetGroups,
		func(v model.ELBTargetGroup) []string {
			// TODO add attached load balancer
			port := "-"
			if v.Port != nil {
				port = strconv.Itoa(int(*v.Port))
			}
			return []string{
				utils.DerefString(v.TargetGroupName, ""),
				port,
				string(v.Protocol),
				utils.AutoCase(string(v.TargetType)),
				utils.DerefString(v.VpcId, ""),
			}
		},
		func(v *TableView[model.ELBTargetGroup]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.TargetGroupArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.TargetGroupArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ELB.Target Groups", func(r Repos, app *Application) Component {
		return NewELBTargetGroups(r.ELB, app)
	})
}
