package main

import (
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ELBLoadBalancers struct {
	*ui.Table
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

func (e ELBLoadBalancers) GetService() string {
	return "ELB"
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
	name, err := e.GetColSelection("NAME")
	if err != nil {
		return
	}
	if arn := e.model[row-1].LoadBalancerArn; arn != nil {
		tagsView := NewELBTags(e.repo, *arn, name, e.app)
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
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
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
		var name, dnsName, lbType, vpcId string
		if v.LoadBalancerName != nil {
			name = *v.LoadBalancerName
		}
		if v.DNSName != nil {
			dnsName = *v.DNSName
		}
		if v.VpcId != nil {
			vpcId = *v.VpcId
		}
		lbType = utils.TitleCase(string(v.Type))
		data = append(data, []string{
			name,
			dnsName,
			lbType,
			vpcId,
		})
	}
	e.SetData(data)
}
