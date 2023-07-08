package main

import (
	"context"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

type Route53HostedZones struct {
	*Table
	r53Client *r53.Client
	app       *Application
}

func NewRoute53HostedZones(client *r53.Client, app *Application) *Route53HostedZones {
	r := &Route53HostedZones{
		Table: NewTable([]string{
			"ID",
			"NAME",
			"RECORDS",
			"VISIBILITY",
			"DESCRIPTION",
		}, 1, 0),
		r53Client: client,
		app:       app,
	}
	r.SetSelectedFunc(r.selectHandler)
	return r
}

func (r Route53HostedZones) GetService() string {
	return "Route 53"
}

func (r Route53HostedZones) GetLabels() []string {
	return []string{"Hosted Zones"}
}

func (r Route53HostedZones) selectHandler(row, col int) {
	hostedZoneId, err := r.GetColSelection("ID")
	if err != nil {
		return
	}
	recordsView := NewRoute53Records(r.r53Client, hostedZoneId)
	r.app.AddAndSwitch(recordsView)
}

func (r Route53HostedZones) tagsHandler() {
	hostedZoneId, err := r.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewRoute53Tags(r.r53Client, r53Types.TagResourceTypeHostedzone, hostedZoneId, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r Route53HostedZones) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r Route53HostedZones) Render() {
	pg := r53.NewListHostedZonesPaginator(
		r.r53Client,
		&r53.ListHostedZonesInput{},
	)
	var hostedZones []r53Types.HostedZone
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		hostedZones = append(hostedZones, out.HostedZones...)
	}

	var data [][]string
	for _, v := range hostedZones {
		var comment, visibility string
		if v.Config != nil {
			if v.Config.Comment != nil {
				comment = *v.Config.Comment
			}
			if v.Config.PrivateZone {
				visibility = "Private"
			} else {
				visibility = "Public"
			}
		}
		splitId := strings.Split(*v.Id, "/")
		data = append(data, []string{
			splitId[len(splitId)-1],
			*v.Name,
			strconv.FormatInt(*v.ResourceRecordSetCount, 10),
			visibility,
			comment,
		})
	}
	r.SetData(data)
}
