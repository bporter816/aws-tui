package main

import (
	"context"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/gdamore/tcell/v2"
)

type Route53HealthChecks struct {
	*ui.Table
	repo      *repo.Route53
	r53Client *r53.Client
	app       *Application
}

func NewRoute53HealthChecks(repo *repo.Route53, r53Client *r53.Client, app *Application) *Route53HealthChecks {
	r := &Route53HealthChecks{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"TYPE",
			"DESCRIPTION",
		}, 1, 0),
		repo:      repo,
		r53Client: r53Client,
		app:       app,
	}
	return r
}

func (r Route53HealthChecks) GetService() string {
	return "Route 53"
}

func (r Route53HealthChecks) GetLabels() []string {
	return []string{"Health Checks"}
}

func (r Route53HealthChecks) tagsHandler() {
	healthCheckId, err := r.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewRoute53Tags(r.repo, r53Types.TagResourceTypeHealthcheck, healthCheckId, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r Route53HealthChecks) GetKeyActions() []KeyAction {
	return []KeyAction{
		KeyAction{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r Route53HealthChecks) Render() {
	pg := r53.NewListHealthChecksPaginator(
		r.r53Client,
		&r53.ListHealthChecksInput{},
	)
	var healthchecks []r53Types.HealthCheck
	for pg.HasMorePages() {
		out, err := pg.NextPage(context.TODO())
		if err != nil {
			panic(err)
		}
		healthchecks = append(healthchecks, out.HealthChecks...)
	}

	var data [][]string
	for _, v := range healthchecks {
		var id, name, checkType, description string
		if v.Id != nil {
			id = *v.Id

			// name comes from the Name tag
			tags, err := r.repo.ListTags(*v.Id, r53Types.TagResourceTypeHealthcheck)
			if err != nil {
				panic(err)
			}
			if n, ok := tags.Get("Name"); ok {
				name = n
			}
		}
		if v.HealthCheckConfig != nil {
			checkType = string(v.HealthCheckConfig.Type)
			description = getHealthCheckDescription(*v.HealthCheckConfig)
		}
		data = append(data, []string{
			id,
			name,
			checkType,
			description,
		})
	}
	r.SetData(data)
}
