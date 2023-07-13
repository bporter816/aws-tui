package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbTypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
)

// TODO support load balancer lists by lb arns
// TODO better manage name/arn
type ELBListeners struct {
	*Table
	elbClient *elb.Client
	app       *Application
	lbArn     string
	lbName    string
	arns      []string
}

func NewELBListeners(elbClient *elb.Client, app *Application, lbArn string, lbName string) *ELBListeners {
	e := &ELBListeners{
		Table: NewTable([]string{
			"PROTOCOL",
			"PORT",
			"RULES",
		}, 1, 0),
		elbClient: elbClient,
		app:       app,
		lbArn:     lbArn,
		lbName:    lbName,
	}
	return e
}

func (e ELBListeners) GetService() string {
	return "ELB"
}

func (e ELBListeners) GetLabels() []string {
	return []string{e.lbName, "Listeners"}
}

func (e ELBListeners) tagsHandler() {
	row, err := e.GetRowSelection()
	if err != nil {
		return
	}
	protocol, err := e.GetColSelection("PROTOCOL")
	if err != nil {
		return
	}
	port, err := e.GetColSelection("PORT")
	if err != nil {
		return
	}
	tagsView := NewELBTags(e.elbClient, ELBResourceTypeListener, e.arns[row-1], protocol+port, e.app)
	e.app.AddAndSwitch(tagsView)
}

func (e ELBListeners) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      e.tagsHandler,
		},
	}
}

func (e *ELBListeners) Render() {
	pg := elb.NewDescribeListenersPaginator(
		e.elbClient,
		&elb.DescribeListenersInput{
			LoadBalancerArn: aws.String(e.lbArn),
		},
	)
	var listeners []elbTypes.Listener
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		listeners = append(listeners, out.Listeners...)
	}

	listenerToRules := make(map[string][]elbTypes.Rule)
	for _, l := range listeners {
		var marker *string
		for {
			out, err := e.elbClient.DescribeRules(
				context.TODO(),
				&elb.DescribeRulesInput{
					ListenerArn: l.ListenerArn,
					Marker:      marker,
				},
			)
			if err != nil {
				panic(err)
			}
			listenerToRules[*l.ListenerArn] = append(listenerToRules[*l.ListenerArn], out.Rules...)
			marker = out.NextMarker
			if marker == nil {
				break
			}
		}
	}

	var data [][]string
	e.arns = make([]string, len(listeners))
	for i, v := range listeners {
		e.arns[i] = *v.ListenerArn
		var protocol, port, rules string
		protocol = string(v.Protocol)
		if v.Port != nil {
			port = strconv.Itoa(int(*v.Port))
		}
		rules = strconv.Itoa(len(listenerToRules[*v.ListenerArn]))
		data = append(data, []string{
			protocol,
			port,
			rules,
		})
	}
	e.SetData(data)
}
