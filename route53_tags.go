package main

import (
	r53Types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
)

type Route53Tags struct {
	*ui.Table
	repo         *repo.Route53
	resourceType r53Types.TagResourceType
	resourceName string
	app          *Application
}

func NewRoute53Tags(repo *repo.Route53, resourceType r53Types.TagResourceType, resourceName string, app *Application) *Route53Tags {
	r := &Route53Tags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:         repo,
		resourceType: resourceType,
		resourceName: resourceName,
		app:          app,
	}
	return r
}

func (r Route53Tags) GetService() string {
	return "Route 53"
}

func (r Route53Tags) GetLabels() []string {
	return []string{r.resourceName, "Tags"}
}

func (r Route53Tags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r Route53Tags) Render() {
	model, err := r.repo.ListTags(r.resourceName, r.resourceType)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		data = append(data, []string{
			v.Key,
			v.Value,
		})
	}
	r.SetData(data)
}
