package internal

import (
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/bporter816/aws-tui/internal/model"
	"github.com/bporter816/aws-tui/internal/repo"
	"github.com/bporter816/aws-tui/internal/utils"
	"github.com/gdamore/tcell/v2"
)

func NewSNSTopics(repo *repo.SNS, app *Application) *TableView[model.SNSTopic] {
	return NewTableView(
		"SNS",
		[]string{"Topics"},
		[]string{"NAME", "TYPE", "PENDING SUBS", "CONFIRMED SUBS", "DELETED SUBS"},
		repo.ListTopics,
		func(v model.SNSTopic) []string {
			a, err := arn.Parse(v.Arn)
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
			return []string{
				utils.GetResourceNameFromArn(a),
				topicType,
				pendingSubs,
				confirmedSubs,
				deletedSubs,
			}
		},
		func(v *TableView[model.SNSTopic]) []KeyAction {
			return []KeyAction{
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 's', tcell.ModNone),
					Description: "Subscriptions",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSNSSubscriptions(repo, item.Arn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone),
					Description: "Access Control Policy",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSNSAccessControlPolicy(repo, item.Arn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone),
					Description: "Delivery Policy",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewSNSDeliveryPolicy(repo, item.Arn, app))
					},
				},
				{
					Key:         tcell.NewEventKey(tcell.KeyRune, 'T', tcell.ModNone),
					Description: "Tags",
					Action: func() {
						item, err := v.Selected()
						if err != nil {
							return
						}
						app.AddAndSwitch(NewTags(repo, v.GetService(), item.Arn, app))
					},
				},
			}
		},
	)
}

func init() {
	register("SNS.Topics", func(r Repos, app *Application) Component {
		return NewSNSTopics(r.SNS, app)
	})
}
