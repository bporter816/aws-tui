package main

import (
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"strings"
)

type SQSAccessPolicy struct {
	*ui.Text
	repo     *repo.SQS
	queueUrl string
	app      *Application
}

func NewSQSAccessPolicy(repo *repo.SQS, queueUrl string, app *Application) *SQSAccessPolicy {
	s := &SQSAccessPolicy{
		Text:     ui.NewText(true, "json"),
		repo:     repo,
		queueUrl: queueUrl,
		app:      app,
	}
	return s
}

func (s SQSAccessPolicy) GetService() string {
	return "SQS"
}

func (s SQSAccessPolicy) GetLabels() []string {
	parts := strings.Split(s.queueUrl, "/")
	if len(parts) > 0 {
		return []string{parts[len(parts)-1], "Access Policy"}
	} else {
		return []string{s.queueUrl, "Access Policy"}
	}
}

func (s SQSAccessPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SQSAccessPolicy) Render() {
	policy, err := s.repo.GetAccessPolicy(s.queueUrl)
	if err != nil {
		panic(err)
	}
	s.SetText(policy)
}
