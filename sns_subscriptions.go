package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type SNSSubscriptions struct {
	*ui.Table
	repo     *repo.SNS
	topicArn string
	app      *Application
}

func NewSNSSubscriptions(repo *repo.SNS, topicArn string, app *Application) *SNSSubscriptions {
	s := &SNSSubscriptions{
		Table: ui.NewTable([]string{
			"ID",
			"PROTOCOL",
			"ENDPOINT",
		}, 1, 0),
		repo:     repo,
		topicArn: topicArn,
		app:      app,
	}
	return s
}

func (s SNSSubscriptions) GetService() string {
	return "SNS"
}

func (s SNSSubscriptions) GetLabels() []string {
	arn, err := arn.Parse(s.topicArn)
	if err != nil {
		panic(err)
	}
	return []string{utils.GetResourceNameFromArn(arn), "Subscriptions"}
}

func (s SNSSubscriptions) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSSubscriptions) Render() {
	model, err := s.repo.ListSubscriptions(s.topicArn)
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		var id, protocol, endpoint string
		if v.SubscriptionArn != nil {
			arn, err := arn.Parse(*v.SubscriptionArn)
			if err != nil {
				panic(err)
			}
			id = utils.GetResourceNameFromArn(arn)
		}
		if v.Protocol != nil {
			protocol = *v.Protocol
		}
		if v.Endpoint != nil {
			endpoint = *v.Endpoint
		}
		data = append(data, []string{
			id,
			protocol,
			endpoint,
		})
	}
	s.SetData(data)
}
