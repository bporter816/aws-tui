package internal

import (
	"context"

	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

type route53HealthCheckRow struct {
	model.Route53HealthCheck
	name string
}

func NewRoute53HealthChecks(repo *repo.Route53, app *Application) *TableView[route53HealthCheckRow] {
	return NewTableView(
		"Route 53",
		[]string{"Health Checks"},
		[]string{"ID", "NAME", "TYPE", "DESCRIPTION"},
		func(ctx context.Context) ([]route53HealthCheckRow, error) {
			checks, err := repo.ListHealthChecks(ctx)
			if err != nil {
				return nil, err
			}
			rows := make([]route53HealthCheckRow, len(checks))
			for i, c := range checks {
				rows[i] = route53HealthCheckRow{Route53HealthCheck: c}
				if c.Id != nil {
					// name comes from the Name tag
					if tags, err := repo.ListTags(ctx, string(r53Types.TagResourceTypeHealthcheck)+":"+*c.Id); err == nil {
						if n, ok := tags.Get("Name"); ok {
							rows[i].name = n
						}
					}
				}
			}
			return rows, nil
		},
		func(v route53HealthCheckRow) []string {
			var id, checkType, description string
			if v.Id != nil {
				id = *v.Id
			}
			if v.HealthCheckConfig != nil {
				checkType = string(v.HealthCheckConfig.Type)
				description = utils.FormatRoute53HealthCheckDescription(*v.HealthCheckConfig)
			}
			return []string{
				id,
				v.name,
				checkType,
				description,
			}
		},
		func(v *TableView[route53HealthCheckRow]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						healthCheckId, err := v.GetColSelection("ID")
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), string(r53Types.TagResourceTypeHealthcheck)+":"+healthCheckId, app))
					},
				},
			}
		},
	)
}

func init() {
	register("Route 53.Health Checks", func(r Repos, app *Application) Component {
		return NewRoute53HealthChecks(r.Route53, app)
	})
}

