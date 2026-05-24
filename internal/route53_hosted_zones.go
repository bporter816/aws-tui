package internal

import (
	"strconv"
	"strings"

	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRoute53HostedZones(repo *repo.Route53, app *Application) *TableView[model.Route53HostedZone] {
	v := NewTableView(
		"Route 53",
		[]string{"Hosted Zones"},
		[]string{"ID", "NAME", "RECORDS", "VISIBILITY", "DESCRIPTION"},
		repo.ListHostedZones,
		func(z model.Route53HostedZone) []string {
			var id, resourceRecordSetCount, visibility, comment string
			if z.Id != nil {
				split := strings.Split(*z.Id, "/")
				id = split[len(split)-1]
			}
			if z.ResourceRecordSetCount != nil {
				resourceRecordSetCount = strconv.FormatInt(*z.ResourceRecordSetCount, 10)
			}
			if z.Config != nil {
				if z.Config.Comment != nil {
					comment = *z.Config.Comment
				}
				if z.Config.PrivateZone {
					visibility = "Private"
				} else {
					visibility = "Public"
				}
			}
			return []string{
				id,
				utils.DerefString(z.Name, ""),
				resourceRecordSetCount,
				visibility,
				comment,
			}
		},
		func(v *TableView[model.Route53HostedZone]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						hostedZoneId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), string(r53Types.TagResourceTypeHostedzone)+":"+hostedZoneId, app))
					},
				},
			}
		},
	)
	v.SetSelectedFunc(func(row, col int) {
		hostedZoneId, err := v.GetColSelection("ID")
		if err != nil {
			return
		}
		recordsView := NewRoute53Records(repo, hostedZoneId, app)
		app.AddAndSwitch(recordsView)
	})
	return v
}
