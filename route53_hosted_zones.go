package main

import (
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
	"strconv"
	"strings"
)

type Route53HostedZones struct {
	*ui.Table
	repo *repo.Route53
	app  *Application
}

func NewRoute53HostedZones(repo *repo.Route53, app *Application) *Route53HostedZones {
	r := &Route53HostedZones{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"RECORDS",
			"VISIBILITY",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
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
	recordsView := NewRoute53Records(r.repo, hostedZoneId, r.app)
	r.app.AddAndSwitch(recordsView)
}

func (r Route53HostedZones) tagsHandler() {
	hostedZoneId, err := r.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewRoute53Tags(r.repo, r53Types.TagResourceTypeHostedzone, hostedZoneId, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r Route53HostedZones) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r Route53HostedZones) Render() {
	model, err := r.repo.ListHostedZones()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, name, resourceRecordSetCount, visibility, comment string
		if v.Id != nil {
			split := strings.Split(*v.Id, "/")
			id = split[len(split)-1]
		}
		if v.Name != nil {
			name = *v.Name
		}
		if v.ResourceRecordSetCount != nil {
			resourceRecordSetCount = strconv.FormatInt(*v.ResourceRecordSetCount, 10)
		}
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
		data = append(data, []string{
			id,
			name,
			resourceRecordSetCount,
			visibility,
			comment,
		})
	}
	r.SetData(data)
}
