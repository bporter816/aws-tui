package internal

import (
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type ELBLoadBalancers struct {
	*ui.Table
	view.ELB
	repo  *repo.ELB
	app   *Application
	model []model.ELBLoadBalancer
}

func NewELBLoadBalancers(repo *repo.ELB, app *Application) *ELBLoadBalancers {
	e := &ELBLoadBalancers{
		Table: ui.NewTable([]string{
			"NAME",
			"DNS NAME",
			"TYPE",
			"VPC",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return e
}

func (e ELBLoadBalancers) GetLabels() []string {
	return []string{"Load Balancers"}
}

func (e ELBLoadBalancers) listenersHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if arn := e.model[row-1].LoadBalancerArn; arn != nil {
		listenersView := NewELBListeners(e.repo, *arn, name, e.app)
		e.app.AddAndSwitch(listenersView)
	}
}

func (e ELBLoadBalancers) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	if arn := e.model[row-1].LoadBalancerArn; arn != nil {
		tagsView := NewTags(e.repo, e.GetService(), *arn, e.app)
		e.app.AddAndSwitch(tagsView)
	}
}

func (e ELBLoadBalancers) GetKeyActions() []KeyAction {
	// TODO add security groups hotkey
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone),
			Description: "Listeners",
			Action:      e.listenersHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBLoadBalancers) Render() {
	model, err := e.repo.ListLoadBalancers()
	if err != nil {
		panic(err)
	}
	e.model = model

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			utils.DerefString(v.LoadBalancerName, ""),
			utils.DerefString(v.DNSName, ""),
			utils.TitleCase(string(v.Type)),
			utils.DerefString(v.VpcId, ""),
		})
	}
	e.SetData(data)
}
