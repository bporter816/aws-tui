package internal

import (
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewRoute53HealthChecks(repo *repo.Route53, app *Application) *TableView[model.Route53HealthCheck] {
	return NewTableView(
		"Route 53",
		[]string{"Health Checks"},
		[]string{"ID", "NAME", "TYPE", "DESCRIPTION"},
		repo.ListHealthChecks,
		func(v model.Route53HealthCheck) []string {
			var id, name, checkType, description string
			if v.Id != nil {
				id = *v.Id

				// name comes from the Name tag
				tags, err := repo.ListTags(string(r53Types.TagResourceTypeHealthcheck) + ":" + *v.Id)
				if err != nil {
					panic(err)
				}
				if n, ok := tags.Get("Name"); ok {
					name = n
				}
			}
			if v.HealthCheckConfig != nil {
				checkType = string(v.HealthCheckConfig.Type)
				description = utils.FormatRoute53HealthCheckDescription(*v.HealthCheckConfig)
			}
			return []string{
				id,
				name,
				checkType,
				description,
			}
		},
		func(v *TableView[model.Route53HealthCheck]) []KeyAction {
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
