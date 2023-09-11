package main

import (
	"context"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type ELBLoadBalancers struct {
	*ui.Table
	elbClient *elb.Client
	app       *Application
	arns      []string
}

func NewELBLoadBalancers(elbClient *elb.Client, app *Application) *ELBLoadBalancers {
	e := &ELBLoadBalancers{
		Table: ui.NewTable([]string{
			"NAME",
			"DNS NAME",
			"TYPE",
			"VPC",
		}, 1, 0),
		elbClient: elbClient,
		app:       app,
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
	listenersView := NewELBListeners(e.elbClient, e.arns[row-1], name, e.app)
	e.app.AddAndSwitch(listenersView)
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
	tagsView := NewELBTags(e.elbClient, ELBResourceTypeLoadBalancer, e.arns[row-1], name, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ELBLoadBalancers) GetKeyActions() []KeyAction {
	// TODO add security groups hotkey
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone),
			Description: "Listeners",
			Action:      e.listenersHandler,
		},
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBLoadBalancers) Render() {
	pg := elb.NewDescribeLoadBalancersPaginator(
		e.elbClient,
		&elb.DescribeLoadBalancersInput{},
	)
	var loadBalancers []elbTypes.LoadBalancer
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		loadBalancers = append(loadBalancers, out.LoadBalancers...)
	}

	var data [][]string
	e.arns = make([]string, len(loadBalancers))
	for i, v := range loadBalancers {
		e.arns[i] = *v.LoadBalancerArn
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
