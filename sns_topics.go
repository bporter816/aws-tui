package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
)

type SNSTopics struct {
	*ui.Table
	repo *repo.SNS
	app  *Application
}

func NewSNSTopics(repo *repo.SNS, app *Application) *SNSTopics {
	s := &SNSTopics{
		Table: ui.NewTable([]string{
			"NAME",
			"TYPE",
			"PENDING SUBS",
			"CONFIRMED SUBS",
			"DELETED SUBS",
		}, 1, 0),
		repo: repo,
		app:  app,
	}
	return s
}

func (s SNSTopics) GetService() string {
	return "SNS"
}

func (s SNSTopics) GetLabels() []string {
	return []string{"Topics"}
}

func (s SNSTopics) GetKeyActions() []KeyAction {
	return []KeyAction{}
}

func (s SNSTopics) Render() {
	model, err := s.repo.ListTopics()
	if err != nil {
		panic(err)
	}

	var data [][]string
	for _, v := range model {
		arn, err := arn.Parse(v.Arn)
		if err != nil {
			panic(err)
		}
		var topicType, pendingSubs, confirmedSubs, deletedSubs string
		if len(v.Attributes) > 0 {
			topicType = "Standard"
			if isFifo, ok := v.Attributes["FifoTopic"]; ok && isFifo == "true" {
				topicType = "FIFO"
			}
			if p, ok := v.Attributes["SubscriptionsPending"]; ok {
				pendingSubs = p
			}
			if c, ok := v.Attributes["SubscriptionsConfirmed"]; ok {
				confirmedSubs = c
			}
			if d, ok := v.Attributes["SubscriptionsDeleted"]; ok {
				deletedSubs = d
			}
		}
		data = append(data, []string{
			utils.GetResourceNameFromArn(arn),
			topicType,
			pendingSubs,
			confirmedSubs,
			deletedSubs,
		})
	}
	s.SetData(data)
}
