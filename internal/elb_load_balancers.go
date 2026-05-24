package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewELBLoadBalancers(repo *repo.ELB, app *Application) *TableView[model.ELBLoadBalancer] {
	return NewTableView(
		"ELB",
		[]string{"Load Balancers"},
		[]string{"NAME", "DNS NAME", "TYPE", "VPC"},
		repo.ListLoadBalancers,
		func(v model.ELBLoadBalancer) []string {
			return []string{
				utils.DerefString(v.LoadBalancerName, ""),
				utils.DerefString(v.DNSName, ""),
				utils.TitleCase(string(v.Type)),
				utils.DerefString(v.VpcId, ""),
			}
		},
		func(v *TableView[model.ELBLoadBalancer]) []KeyAction {
			// TODO add security groups hotkey
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone),
					Description: "Listeners",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.LoadBalancerArn == nil {
							return
						}
						name, err := v.GetColSelection("NAME")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewELBListeners(repo, *item.LoadBalancerArn, name, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil || item.LoadBalancerArn == nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), *item.LoadBalancerArn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("ELB.Load Balancers", func(r Repos, app *Application) Component {
		return NewELBLoadBalancers(r.ELB, app)
	})
}
