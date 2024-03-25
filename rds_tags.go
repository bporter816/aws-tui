package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type RDSTags struct {
	*ui.Table
	view.RDS
	repo       *repo.RDS
	resourceId string
	app        *Application
}

func NewRDSTags(repo *repo.RDS, resourceId string, app *Application) *RDSTags {
	r := &RDSTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:       repo,
		resourceId: resourceId,
		app:        app,
	}
	return r
}

func (r RDSTags) GetLabels() []string {
	arn, err := arn.Parse(r.resourceId)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (r RDSTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (r RDSTags) Render() {
	model, err := r.repo.ListTags(r.resourceId)
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
