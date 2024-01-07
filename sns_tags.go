package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/bporter816/aws-tui/view"
)

type SNSTags struct {
	*ui.Table
	view.SNS
	repo     *repo.SNS
	topicArn string
	app      *Application
}

func NewSNSTags(repo *repo.SNS, topicArn string, app *Application) *SNSTags {
	s := &SNSTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:     repo,
		topicArn: topicArn,
		app:      app,
	}
	return s
}

func (s SNSTags) GetLabels() []string {
	arn, err := arn.Parse(s.topicArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Tags"}
}

func (s SNSTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSTags) Render() {
	model, err := s.repo.ListTags(s.topicArn)
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
	s.SetData(data)
}
