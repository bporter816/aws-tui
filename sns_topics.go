package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/model"
	"github.com/bporter816/aws-tui/repo"
	"github.com/bporter816/aws-tui/ui"
	"github.com/bporter816/aws-tui/utils"
	"github.com/gdamore/tcell/v2"
)

type SNSTopics struct {
	*ui.Table
	repo  *repo.SNS
	app   *Application
	model []model.SNSTopic
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

func (s SNSTopics) subscriptionsHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	subscriptionsView := NewSNSSubscriptions(s.repo, s.model[row-1].Arn, s.app)
	s.app.AddAndSwitch(subscriptionsView)
}

func (s SNSTopics) accessControlPolicyHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	accessControlPolicyView := NewSNSAccessControlPolicy(s.repo, s.model[row-1].Arn, s.app)
	s.app.AddAndSwitch(accessControlPolicyView)
}

func (s SNSTopics) deliveryPolicyHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	deliveryPolicyView := NewSNSDeliveryPolicy(s.repo, s.model[row-1].Arn, s.app)
	s.app.AddAndSwitch(deliveryPolicyView)
}

func (s SNSTopics) tagsHandler() {
	row, err := s.GetRowSelection()
	if err != nil {
		return
	}
	tagsView := NewSNSTags(s.repo, s.model[row-1].Arn, s.app)
	s.app.AddAndSwitch(tagsView)
}

func (s SNSTopics) GetKeyActions() []KeyAction {
	return []KeyAction{
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
			Description: "Subscriptions",
			Action:      s.subscriptionsHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
			Description: "Access Control Policy",
			Action:      s.accessControlPolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
			Description: "Delivery Policy",
			Action:      s.deliveryPolicyHandler,
		},
		{
			Key:         tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone),
			Description: "Tags",
			Action:      s.tagsHandler,
		},
	}
}

func (s *SNSTopics) Render() {
	model, err := s.repo.ListTopics()
	if err != nil {
		panic(err)
	}
	s.model = model

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
