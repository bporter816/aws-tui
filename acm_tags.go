package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type ACMTags struct {
	*ui.Table
	view.ACM
	repo           *repo.ACM
	certificateArn arn.ARN
	app            *Application
}

func NewACMTags(repo *repo.ACM, certificateArn string, app *Application) *ACMTags {
	arn, err := arn.Parse(certificateArn)
	if err != nil {
		panic(err)
	}
	a := &ACMTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:           repo,
		certificateArn: arn,
		app:            app,
	}
	return a
}

func (a ACMTags) GetLabels() []string {
	return []string{utils.GetResourceNameFromArn(a.certificateArn), "Tags"}
}

func (a ACMTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (a ACMTags) Render() {
	model, err := a.repo.ListTags(a.certificateArn)
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
	a.SetData(data)
}
