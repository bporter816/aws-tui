package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	r53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53HealthChecks struct {
	*Table
	r53Client *r53.Client
	app       *Application
}

func NewRoute53HealthChecks(r53Client *r53.Client, app *Application) *Route53HealthChecks {
	r := &Route53HealthChecks{
		Table: NewTable([]string{
			"ID",
			"NAME",
			"TYPE",
			"DESCRIPTION",
		}, 1, 0),
		r53Client: r53Client,
		app:       app,
	}
	return r
}

func (r Route53HealthChecks) GetName() string {
	return "Route 53 | Health Checks"
}

func (r Route53HealthChecks) GetKeyActions() []KeyAction {
	return []KeyAction{}
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
		// name comes from the Name tag
		out, err := r.r53Client.ListTagsForResource(
			context.TODO(),
			&r53.ListTagsForResourceInput{
				ResourceId:   aws.String(*v.Id),
				ResourceType: r53Types.TagResourceTypeHealthcheck,
			},
		)
		if err != nil {
			panic(err)
		}

		var name, checkType, description string
		if out.ResourceTagSet != nil {
			name, _ = getTag(out.ResourceTagSet.Tags, "Name")
		}
		if v.HealthCheckConfig != nil {
			checkType = string(v.HealthCheckConfig.Type)
			description = getHealthCheckDescription(*v.HealthCheckConfig)
		}
		data = append(data, []string{
			*v.Id,
			name,
			checkType,
			description,
		})
	}
	r.SetData(data)
}
