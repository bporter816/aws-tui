package internal

import (
	"strconv"
	"strings"

	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/ui"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/bporter816/aws-tui/internal/view"
	"github.com/gdamore/tcell/v2"
)

type Route53HostedZones struct {
	*ui.Table
	view.Route53
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
	tagsView := NewTags(r.repo, r.GetService(), string(r53Types.TagResourceTypeHostedzone)+":"+hostedZoneId, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r Route53HostedZones) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
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
		var id, resourceRecordSetCount, visibility, comment string
		if v.Id != nil {
			split := strings.Split(*v.Id, "/")
			id = split[len(split)-1]
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
			utils.DerefString(v.Name, ""),
			resourceRecordSetCount,
			visibility,
			comment,
		})
	}
	r.SetData(data)
}
