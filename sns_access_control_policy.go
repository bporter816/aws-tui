package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type SNSAccessControlPolicy struct {
	*ui.Text
	repo     *repo.SNS
	topicArn string
	app      *Application
}

func NewSNSAccessControlPolicy(repo *repo.SNS, topicArn string, app *Application) *SNSAccessControlPolicy {
	s := &SNSAccessControlPolicy{
		Text:     ui.NewText(true, "json"),
		repo:     repo,
		topicArn: topicArn,
		app:      app,
	}
	return s
}

func (s SNSAccessControlPolicy) GetService() string {
	return "SNS"
}

func (s SNSAccessControlPolicy) GetLabels() []string {
	arn, err := arn.Parse(s.topicArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Access Control Policy"}
}

func (s SNSAccessControlPolicy) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSAccessControlPolicy) Render() {
	policy, err := s.repo.GetAccessControlPolicy(s.topicArn)
	if err != nil {
		panic(err)
	}
	s.SetText(policy)
}
