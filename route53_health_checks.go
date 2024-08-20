package main

import (
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
	"github.com/gdamore/tcell/v2"
)

type Route53HealthChecks struct {
	*ui.Table
	view.Route53
	repo *repo.Route53
	app  *Application
}

func NewRoute53HealthChecks(repo *repo.Route53, app *Application) *Route53HealthChecks {
	r := &Route53HealthChecks{
		Table: ui.NewTable([]string{
			"ID",
			"NAME",
			"TYPE",
			"DESCRIPTION",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return r
}

func (r Route53HealthChecks) GetLabels() []string {
	return []string{"Health Checks"}
}

func (r Route53HealthChecks) tagsHandler() {
	healthCheckId, err := r.GetColSelection("ID")
	if err != nil {
		return
	}
	tagsView := NewTags(r.repo, r.GetService(), string(r53Types.TagResourceTypeHealthcheck)+":"+healthCheckId, r.app)
	r.app.AddAndSwitch(tagsView)
}

func (r Route53HealthChecks) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
			Description: "Tags",
			Action:      r.tagsHandler,
		},
	}
}

func (r Route53HealthChecks) Render() {
	model, err := r.repo.ListHealthChecks()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, name, checkType, description string
		if v.Id != nil {
			id = *v.Id

			// name comes from the Name tag
			tags, err := r.repo.ListTags(string(r53Types.TagResourceTypeHealthcheck) + ":" + *v.Id)
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
		data = append(data, []string{
			id,
			name,
			checkType,
			description,
		})
	}
	r.SetData(data)
}
