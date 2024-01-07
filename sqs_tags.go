package main

import (
	"strings"

	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/view"
)

type SQSTags struct {
	*ui.Table
	view.SQS
	repo     *repo.SQS
	queueUrl string
	app      *Application
}

func NewSQSTags(repo *repo.SQS, queueUrl string, app *Application) *SQSTags {
	s := &SQSTags{
		Table: ui.NewTable([]string{
			"KEY",
			"VALUE",
		}, 1, 0),
		repo:     repo,
		queueUrl: queueUrl,
		app:      app,
	}
	return s
}

func (s SQSTags) GetLabels() []string {
	parts := strings.Split(s.queueUrl, "/")
	if len(parts) > 0 {
		return []string{parts[len(parts)-1], "Tags"}
	} else {
		return []string{s.queueUrl, "Tags"}
	}
}

func (s SQSTags) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SQSTags) Render() {
	model, err := s.repo.ListTags(s.queueUrl)
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
